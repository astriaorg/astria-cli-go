package keys

import (
	"github.com/99designs/keyring"
	log "github.com/sirupsen/logrus"
)

const service = "astria-go"

// StoreKeyring stores a secret in the keyring for a user.
func StoreKeyring(key string, secret string) error {
	ring, err := keyring.Open(keyring.Config{
		ServiceName: service,
	})
	if err != nil {
		log.WithError(err).Error("error opening keyring service")
		return err
	}

	err = ring.Set(keyring.Item{
		Key:   key,
		Data:  []byte(secret),
		Label: "Astria Sequencer account private key",
	})
	if err != nil {
		log.WithError(err).Error("error setting key")
		return err
	}

	log.Debug("stored secret using keyring service")
	return nil
}

// GetKeyring gets a secret from the keyring for a user.
func GetKeyring(key string) (string, error) {
	ring, err := keyring.Open(keyring.Config{
		ServiceName: service,
	})
	if err != nil {
		log.WithError(err).Error("error opening keyring service")
		return "", err
	}

	log.Debugf("retrieving secret for %s", key)
	item, err := ring.Get(key)
	if err != nil {
		log.WithError(err).Error("error getting secret from keyring")
		return "", nil
	}

	return string(item.Data), nil
}
