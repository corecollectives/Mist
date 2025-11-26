package utils

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// DeploymentLogger provides structured logging for deployment operations
type DeploymentLogger struct {
	logger     zerolog.Logger
	deployID   int64
	appID      int64
	commitHash string
}

// NewDeploymentLogger creates a new logger instance for a deployment
func NewDeploymentLogger(depID, appID int64, commit string) *DeploymentLogger {
	return &DeploymentLogger{
		logger: log.With().
			Int64("deployment_id", depID).
			Int64("app_id", appID).
			Str("commit_hash", commit).
			Logger(),
		deployID:   depID,
		appID:      appID,
		commitHash: commit,
	}
}

// Info logs informational messages
func (dl *DeploymentLogger) Info(message string) {
	dl.logger.Info().Msg(message)
}

// InfoWithFields logs informational messages with additional fields
func (dl *DeploymentLogger) InfoWithFields(message string, fields map[string]interface{}) {
	event := dl.logger.Info()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(message)
}

// Error logs error messages
func (dl *DeploymentLogger) Error(err error, message string) {
	dl.logger.Error().Err(err).Msg(message)
}

// ErrorWithFields logs error messages with additional fields
func (dl *DeploymentLogger) ErrorWithFields(err error, message string, fields map[string]interface{}) {
	event := dl.logger.Error().Err(err)
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(message)
}

// Warn logs warning messages
func (dl *DeploymentLogger) Warn(message string) {
	dl.logger.Warn().Msg(message)
}

// WarnWithFields logs warning messages with additional fields
func (dl *DeploymentLogger) WarnWithFields(message string, fields map[string]interface{}) {
	event := dl.logger.Warn()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(message)
}

// Debug logs debug messages
func (dl *DeploymentLogger) Debug(message string) {
	dl.logger.Debug().Msg(message)
}

// InitLogger initializes the global logger configuration
func InitLogger() {
	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
}
