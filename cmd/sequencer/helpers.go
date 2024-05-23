package sequencer

import (
	"fmt"

	"github.com/astria/astria-cli-go/internal/keys"
	"github.com/astria/astria-cli-go/internal/sequencer"
	"github.com/pterm/pterm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

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
	account := sequencer.NewAccountFromPrivKey(privkey)
	return account.PrivateKeyString(), nil
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
