package devtools

import (
	"fmt"
	"regexp"
)

func IsInstanceNameValid(instance string) error {
	pattern := `^[a-z]+[a-z0-9]*(-[a-z0-9]+)*$`
	matched, err := regexp.MatchString(pattern, instance)
	if err != nil || !matched {
		return fmt.Errorf("Invalid instance name: '%s' Instance names must be lowercase, alphanumeric, and may contain dashes. It can't begin or end with dash. No repeating dashes.", instance)
	}
	return nil
}
