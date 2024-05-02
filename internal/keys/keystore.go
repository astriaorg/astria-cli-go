package keys

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/scrypt"
)

// EncryptedKeyStore defines the structure of the encrypted keystore.
type EncryptedKeyStore struct {
	Address string `json:"address"`
	Crypto  struct {
		Ciphertext   string `json:"ciphertext"`
		Cipherparams struct {
			IV string `json:"iv"`
		} `json:"cipherparams"`
		Cipher    string `json:"cipher"`
		Kdf       string `json:"kdf"`
		Kdfparams struct {
			Dklen int    `json:"dklen"`
			Salt  string `json:"salt"`
			N     int    `json:"n"`
			R     int    `json:"r"`
			P     int    `json:"p"`
		} `json:"kdfparams"`
		MAC string `json:"mac"`
	} `json:"crypto"`
	Version int `json:"version"`
}

// NewEncryptedKeyStore creates a new encrypted keystore using the provided password and private key.
func NewEncryptedKeyStore(password string, address string, priv ed25519.PrivateKey) (*EncryptedKeyStore, error) {
	// generate a salt for the scrypt key derivation function
	salt := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		log.WithError(err).Error("Error reading in random data for salt")
		return nil, err
	}

	// derive a key using scrypt
	derivedKey, err := scrypt.Key([]byte(password), salt, 16384, 8, 1, 32)
	if err != nil {
		log.WithError(err).Error("Error deriving key")
		return nil, err
	}

	// prepare the block cypher
	block, err := aes.NewCipher(derivedKey)
	if err != nil {
		log.WithError(err).Error("Error creating new cipher")
		return nil, err
	}

	// generate an initialization vector. 12 is correct size for gcm
	iv := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.WithError(err).Error("Error reading in random data for initialization vector")
		return nil, err
	}

	// create a new GCM block cipher
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		log.WithError(err).Error("Error creating new GCM block cipher")
		return nil, err
	}

	// encrypt the private key
	ciphertext := aesgcm.Seal(nil, iv, priv, nil)

	// create a MAC for the ciphertext using a slice of the derived key
	mac, err := newHMAC(derivedKey[:16], ciphertext)
	if err != nil {
		log.WithError(err).Error("Error creating new HMAC")
		return nil, err
	}

	// fill the keystore structure.
	keystore := EncryptedKeyStore{
		Address: address,
		Version: 3,
	}
	keystore.Crypto.Cipher = "aes-256-gcm"
	keystore.Crypto.Cipherparams.IV = fmt.Sprintf("%x", iv)
	keystore.Crypto.Ciphertext = fmt.Sprintf("%x", ciphertext)
	keystore.Crypto.Kdf = "scrypt"
	keystore.Crypto.Kdfparams.Dklen = 32
	keystore.Crypto.Kdfparams.Salt = fmt.Sprintf("%x", salt)
	keystore.Crypto.Kdfparams.N = 16384
	keystore.Crypto.Kdfparams.R = 8
	keystore.Crypto.Kdfparams.P = 1
	keystore.Crypto.MAC = fmt.Sprintf("%x", mac)

	return &keystore, nil
}

// DecryptPrivateKey decrypts the private key stored in the keystore using the provided password.
func DecryptPrivateKey(keystore *EncryptedKeyStore, password string) (ed25519.PrivateKey, error) {
	// Decode the salt, IV, and ciphertext from the keystore.
	salt, err := hex.DecodeString(keystore.Crypto.Kdfparams.Salt)
	if err != nil {
		return nil, err
	}
	iv, err := hex.DecodeString(keystore.Crypto.Cipherparams.IV)
	if err != nil {
		return nil, err
	}
	ciphertext, err := hex.DecodeString(keystore.Crypto.Ciphertext)
	if err != nil {
		return nil, err
	}

	// derive the same key from the password and salt
	key, err := scrypt.Key([]byte(password), salt, keystore.Crypto.Kdfparams.N, keystore.Crypto.Kdfparams.R, keystore.Crypto.Kdfparams.P, 32)
	if err != nil {
		return nil, err
	}

	// initialize the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// decrypt the ciphertext
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	privKey, err := aesgcm.Open(nil, iv, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return privKey, nil
}

// newHMAC generates a HMAC for the given key and data.
func newHMAC(key, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.WithError(err).Error("Error creating new cipher")
		return nil, err
	}

	mac := cipher.NewCBCEncrypter(block, key)
	macOutput := make([]byte, len(data))
	mac.CryptBlocks(macOutput, data)
	return macOutput, nil
}

func SaveKeystoreToFile(keydir string, keystore *EncryptedKeyStore) (string, error) {
	bytes, err := json.MarshalIndent(keystore, "", "  ")
	if err != nil {
		log.WithError(err).Error("Cannot marshal Keystore")
		return "", err
	}

	timestamp := time.Now().Format(time.RFC3339)
	filename := fmt.Sprintf("UTC--%s--%s", timestamp, keystore.Address)
	fullpath := filepath.Join(keydir, filename)

	err = os.WriteFile(fullpath, bytes, 0644)
	if err != nil {
		log.WithError(err).Error("Cannot write file")
		return "", err
	}

	return fullpath, nil
}

// DecryptKeyfile decrypts the private key stored in the keystore using the provided password.
// It reads the content of the keyfile, unmarshalls it into an EncryptedKeyStore struct,
// and then calls DecryptPrivateKey to decrypt the private key using the provided password.
// If successful, it returns the decrypted private key. Otherwise, it returns an error.
func DecryptKeyfile(keyfile string, password string) (ed25519.PrivateKey, error) {
	jsonBytes, err := os.ReadFile(keyfile)
	if err != nil {
		return ed25519.PrivateKey{}, err
	}
	ks := &EncryptedKeyStore{}
	err = json.Unmarshal(jsonBytes, ks)
	if err != nil {
		return ed25519.PrivateKey{}, err
	}

	key, err := DecryptPrivateKey(ks, password)
	if err != nil {
		return ed25519.PrivateKey{}, err
	}

	return key, nil
}

// ResolveKeyfilePath resolves the path to a keyfile in the given directory.
// If the directory itself is a keyfile, it returns the absolute path.
// If the keyfile is not found in the directory, an error is returned.
func ResolveKeyfilePath(keydir string) (string, error) {
	keydir, _ = filepath.Abs(keydir)
	fileInfo, err := os.Stat(keydir)
	if err != nil {
		return "", err
	}
	if !fileInfo.IsDir() {
		return keydir, nil
	}

	files, _ := os.ReadDir(keydir)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if strings.HasPrefix(file.Name(), "UTC--") {
			return filepath.Join(keydir, file.Name()), nil
		}
	}

	return "", fmt.Errorf("keyfile is not in %s", keydir)
}
