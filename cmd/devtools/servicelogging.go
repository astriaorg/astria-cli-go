package devtools

import log "github.com/sirupsen/logrus"

func ValidateServiceLogLevelOrPanic(logLevel string) {
	switch logLevel {
	case "debug", "info", "error":
		return
	default:
		log.WithField("service-log-level", logLevel).Fatal("Invalid service log level. Must be one of: 'debug', 'info', 'error'")
		panic("Invalid service log level")
	}

}
