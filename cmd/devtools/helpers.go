package devtools

import (
	"fmt"
	"regexp"

	log "github.com/sirupsen/logrus"
)

func IsInstanceNameValidOrPanic(instance string) {
	re, err := regexp.Compile(`^[a-z]+[a-z0-9]*(-[a-z0-9]+)*$`)
	if err != nil {
		log.WithError(err).Error("Error compiling regex")
		panic(err)
	}
	if !re.MatchString(instance) {
		log.Errorf("Invalid instance name: %s", instance)
		err := fmt.Errorf(`
Invalid instance name: '%s'. Instance names must be lowercase, alphanumeric, 
and may contain dashes. It can't begin or end with a dash. No repeating dashes.
`, instance)
		panic(err)
	}
}
