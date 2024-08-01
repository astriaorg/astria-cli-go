package cmd

import (
	"os"
	"reflect"

	log "github.com/sirupsen/logrus"
)

// CreateDirOrPanic creates a directory with the given name with 0755 permissions.
// If the directory can't be created, it will panic.
func CreateDirOrPanic(dirName string) {
	err := os.MkdirAll(dirName, 0755)
	if err != nil {
		log.WithError(err).Error("Error creating data directory")
		panic(err)
	}
}

// GetFieldValueByTag gets a field's value from a struct by the specified tagName.
func GetFieldValueByTag(obj interface{}, tagName, tagValue string) (reflect.Value, bool) {
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get(tagName)
		if tag == tagValue {
			return val.Field(i), true
		}
	}
	return reflect.Value{}, false
}
