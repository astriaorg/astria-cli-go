package sequencer

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"math/big"
	"regexp"
	"strings"

	primproto "buf.build/gen/go/astria/primitives/protocolbuffers/go/astria/primitive/v1"
	"github.com/astriaorg/astria-cli-go/modules/cli/internal/keys"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
)

// AddPortToURL adds a port to a URL if it doesn't already have one.
// The port is needed for use with the Sequencer Client.
func AddPortToURL(url string) string {
	// check if the URL already has a port
	matched, err := regexp.MatchString(`:\d+$`, url)
	if err != nil {
		log.WithError(err).Error("Error matching string")
		return url
	}
	if matched {
		log.Debug("Port already present in URL: ", url)
		return url
	}
	if strings.Contains(url, "http:") {
		log.Debug("http url detected without a port. Adding port :80 to url: ", url)
		return url + ":80"
	}
	if strings.Contains(url, "https:") {
		log.Debug("https url detected without a port. Adding port :443 to url: ", url)
		return url + ":443"
	}
	return url
}

// PrivateKeyFromText converts a string representation of a private key to an ed25519.PrivateKey.
// It decodes the private key from hex string format and creates a new ed25519.PrivateKey.
func PrivateKeyFromText(privkey string) (ed25519.PrivateKey, error) {
	privKeyBytes, err := hex.DecodeString(privkey)
	if err != nil {
		return nil, err
	}
	from := ed25519.NewKeyFromSeed(privKeyBytes)
	return from, nil
}

// AddressFromText converts a bech32m string representation of an address to an
// Address protobuf. No validation is done on the input string.
func AddressFromText(addr string) *primproto.Address {
	return &primproto.Address{
		Bech32M: addr,
	}
}

// convertToUint128 converts a string to an Uint128 protobuf
func convertToUint128(numStr string) (*primproto.Uint128, error) {
	bigInt := new(big.Int)

	// convert the string to a big.Int
	_, ok := bigInt.SetString(numStr, 10)
	if !ok {
		return nil, fmt.Errorf("failed to convert string to big.Int")
	}

	// check if the number is negative or overflows Uint128
	if bigInt.Sign() < 0 {
		return nil, fmt.Errorf("negative number not allowed")
	} else if bigInt.BitLen() > 128 {
		return nil, fmt.Errorf("value overflows Uint128")
	}

	// split the big.Int into two uint64s
	// convert the big.Int to uint64, which will drop the higher 64 bits
	lo := bigInt.Uint64()
	// shift the big.Int to the right by 64 bits and convert to uint64
	hi := bigInt.Rsh(bigInt, 64).Uint64()
	uint128 := &primproto.Uint128{
		Lo: lo,
		Hi: hi,
	}

	return uint128, nil
}

// strip0xPrefix removes the 0x prefix from a string if present.
func strip0xPrefix(s string) string {
	return strings.TrimPrefix(s, "0x")
}

// PublicKeyFromText converts a hexadecimal string representation of a public
// key to an ed25519.PublicKey. If the input string is not a valid hexadecimal
// string, an error will be returned.
func PublicKeyFromText(addr string) (ed25519.PublicKey, error) {
	addr = strip0xPrefix(addr)
	bytes, err := hex.DecodeString(addr)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// GetPrivateKeyFromFlags retrieves the private key from the command flags.
// If the 'privkey' flag is set, it returns the value of that flag.
// If the 'keyring-address' flag is set, it calls the 'PrivateKeyFromKeyringAddress' function
// to retrieve the private key from the keyring.
// If the 'keyfile' flag is set, it calls the 'PrivateKeyFromKeyfile' function
// to retrieve the private key from the keyfile.
// If none of the flags are set or if the value of 'keyfile' is empty, it returns an error.
// NOTE - requires the flags `keyfile`, `keyring-address`, and `privkey` along with `MarkFlagsOneRequired` and `MarkFlagsMutuallyExclusive`
func GetPrivateKeyFromFlags(c *cobra.Command) (string, error) {
	keyfile := c.Flag("keyfile").Value.String()
	keyringAddress := c.Flag("keyring-address").Value.String()
	priv := c.Flag("privkey").Value.String()

	// NOTE - this isn't very secure but we still support it
	if priv != "" {
		return priv, nil
	}

	// NOTE - this should trigger user's os keyring password prompt
	if keyringAddress != "" {
		return PrivateKeyFromKeyringAddress(keyringAddress)
	}

	if keyfile != "" {
		return PrivateKeyFromKeyfile(keyfile)
	}

	return "", fmt.Errorf("no private key specified")
}

// PrivateKeyFromKeyfile retrieves the private key from the specified keyfile.
func PrivateKeyFromKeyfile(keyfile string) (string, error) {
	kf, err := keys.ResolveKeyfilePath(keyfile)
	if err != nil {
		return "", err
	}

	pwIn := pterm.DefaultInteractiveTextInput.WithMask("*")
	pw, _ := pwIn.Show("Account password:")

	privkey, err := keys.DecryptKeyfile(kf, pw)
	if err != nil {
		log.WithError(err).Error("Error decrypting keyfile")
		return "", err
	}
	return hex.EncodeToString(privkey[:32]), nil
}

// PrivateKeyFromKeyringAddress retrieves the private key from the keyring for a given keyring address.
func PrivateKeyFromKeyringAddress(keyringAddress string) (string, error) {
	key, err := keys.GetKeyring(keyringAddress)
	if err != nil {
		log.WithError(err).Error("Error getting private key from keyring")
		return "", err
	}
	return key, nil
}
