package models

import "time"

type LogSource string

const (
	LogSourceApp    LogSource = "app"
	LogSourceSystem LogSource = "system"
)

type LogLevel string

const (
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
	LogLevelDebug LogLevel = "debug"
)

type Logs struct {
	ID        int64
	Source    LogSource
	SourceID  *int64
	Message   string
	Level     LogLevel
	CreatedAt time.Time
}
