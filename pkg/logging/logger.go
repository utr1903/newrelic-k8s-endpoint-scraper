package config

import (
	"os"

	"github.com/sirupsen/logrus"
)

const (
	// config
	CONFIG__LICENSE_KEY_IS_NOT_PROVIDED               = "license key is not provided! define config.data.newrelic.licenseKey in your helm deployment"
	CONFIG__ACCOUNT_ID_IS_NOT_PROVIDED                = "account id not provided! Define config.data.newrelic.accountId in your helm deployment"
	CONFIG__CONFIG_PATH_IS_NOT_DEFINED                = "config path is not defined"
	CONFIG__CONFIG_FILE_COULD_NOT_BE_READ             = "config file could not be read"
	CONFIG__CONFIG_FILE_COULD_NOT_BE_PARSED_INTO_YAML = "config file could not be parsed into yaml format"
	CONFIG__NO_ENDPOINT_IS_DEFINED                    = "no endpoint is defined"
	CONFIG__ENDPOINT_INFO_IS_MISSING                  = "check your endpoint definitions! type, name and url must be defined"
	CONFIG__ENDPOINT_TYPE_IS_NOT_SUPPORTED            = "only the following types are supported: kvp"

	// scrape
	SCRAPE__HTTP_REQUEST_COULD_NOT_BE_CREATED = "http request could not be created"
	SCRAPE__HTTP_REQUEST_HAS_FAILED           = "http request has failed"
	SCRAPE__ENDPOINT_RETURNED_NOT_OK_STATUS   = "http request has returned not OK status"
	SCRAPE__RESPONSE_BODY_COULD_NOT_BE_PARSED = "response body could not be parsed"

	// forward
	FORWARD__PAYLOAD_COULD_NOT_BE_CREATED      = "payload could not be created"
	FORWARD__HTTP_REQUEST_COULD_NOT_BE_CREATED = "http request could not be created"
	FORWARD__HTTP_REQUEST_HAS_FAILED           = "http request has failed"
	FORWARD__NEW_RELIC_RETURNED_NOT_OK_STATUS  = "http request has returned not OK status"
)

type Logger struct {
	log *logrus.Logger
}

func NewLogger(
	logLevel string,
) *Logger {
	l := logrus.New()
	l.Out = os.Stdout
	l.Formatter = &logrus.JSONFormatter{}

	switch logLevel {
	case "DEBUG":
		l.Level = logrus.DebugLevel
	default:
		l.Level = logrus.ErrorLevel
	}

	return &Logger{
		log: l,
	}
}

func (l *Logger) Log(
	lvl logrus.Level,
	msg string,
) {

	fields := logrus.Fields{
		"instrumentation.provider": "newrelic-kubernetes-endpoint-scraper",
	}

	switch lvl {
	case logrus.ErrorLevel:
		l.log.WithFields(fields).Error(msg)
	default:
		l.log.WithFields(fields).Debug(msg)
	}
}

func (l *Logger) LogWithFields(
	lvl logrus.Level,
	msg string,
	attributes map[string]string,
) {

	fields := logrus.Fields{
		"instrumentation.provider": "newrelic-kubernetes-endpoint-scraper",
	}

	for key, val := range attributes {
		fields[key] = val
	}

	switch lvl {
	case logrus.ErrorLevel:
		l.log.WithFields(fields).Error(msg)
	default:
		l.log.WithFields(fields).Debug(msg)
	}
}
