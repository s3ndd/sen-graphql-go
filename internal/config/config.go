package config

import (
	"time"

	"github.com/s3ndd/sen-go/client"
	"github.com/s3ndd/sen-go/config"
	"github.com/s3ndd/sen-go/log"
)

// Logger builds the log config from the config file and envs.
func Logger() *log.Config {
	return &log.Config{
		LogLevel:               config.String("LOGGER_LEVEL", string(log.LevelInfo)),
		LogFormat:              config.String("LOGGER_FORMAT", "json"),
		Program:                config.String("SOURCE_PROGRAM", "sen-monitoring"),
		Env:                    config.String("ENV", "dev"),
		SentryLogLevel:         config.String("SENTRY_LOGGER_LEVEL", string(log.LevelError)),
		SentryDsn:              config.String("SENTRY_DSN", ""),
		SentryDebug:            config.Bool("SENTRY_DEBUG", true),
		SentryAttachStacktrace: config.Bool("SENTRY_ATTACH_STACK_TRACE", true),
		SentryFlushTimeout:     config.Duration("SENTRY_FLUSH_TIMEOUT", 3*time.Second),
	}
}

// HTTPClient returns a http client config for integration services
func HTTPClient() *client.Config {
	return &client.Config{
		EndpointURL:             config.String("API_GATEWAY_BASEURL", "localhost:8080"),
		UseSecure:               config.Bool("HTTP_CLIENT_USE_SECURE", true),
		IgnoreCertificateErrors: true,
		Timeout:                 config.Duration("HTTP_CLIENT_TIMEOUT", 5*time.Second),
	}
}

// ApiKey retrieves and returns the value of the api key.
func ApiKey(key string) string {
	return config.String(key, "")
}
