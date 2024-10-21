package mysql

import (
	"context"
	"errors"
	"fmt"
	sklogger "github.com/sk-pkg/logger"
	"go.uber.org/zap"
	gormlogger "gorm.io/gorm/logger"
	"time"
)

const (
	// defaultLogLevel is the default log level for the logger.
	defaultLogLevel = gormlogger.Warn
	// defaultIgnoreRecordNotFoundError determines whether to ignore "record not found" errors by default.
	defaultIgnoreRecordNotFoundError = true
	// defaultSlowThreshold is the default duration threshold for slow query logging.
	defaultSlowThreshold = 200 * time.Millisecond
)

// LoggerOption is a function type used to configure the logger.
type LoggerOption func(*logger)

// WithLevel returns a LoggerOption that sets the log level for the logger.
//
// Parameters:
//   - level: A string representing the desired log level ("debug", "info", "warn", or "silent").
//
// Returns:
//   - A LoggerOption function that sets the log level when applied.
//
// Example:
//
//	logger := NewLog(manager, WithLevel("info"))
func WithLevel(level string) LoggerOption {
	var logLevel gormlogger.LogLevel
	switch level {
	case "warn":
		logLevel = gormlogger.Warn
	case "error":
		logLevel = gormlogger.Error
	case "silent":
		logLevel = gormlogger.Silent
	default:
		logLevel = defaultLogLevel // Default level if input is unrecognized
	}
	return func(l *logger) {
		l.logLevel = logLevel
	}
}

// WithIgnoreRecordNotFoundError returns a LoggerOption that sets whether to ignore "record not found" errors.
//
// Parameters:
//   - ignore: A boolean indicating whether to ignore "record not found" errors.
//
// Returns:
//   - A LoggerOption function that sets the ignore flag when applied.
//
// Example:
//
//	logger := NewLog(manager, WithIgnoreRecordNotFoundError(true))
func WithIgnoreRecordNotFoundError(ignore bool) LoggerOption {
	return func(l *logger) {
		l.ignoreRecordNotFoundError = ignore
	}
}

// WithSlowThreshold returns a LoggerOption that sets the threshold for slow query logging.
//
// Parameters:
//   - threshold: A time.Duration representing the slow query threshold.
//
// Returns:
//   - A LoggerOption function that sets the slow threshold when applied.
//
// Example:
//
//	logger := NewLog(manager, WithSlowThreshold(500 * time.Millisecond))
func WithSlowThreshold(threshold time.Duration) LoggerOption {
	return func(l *logger) {
		l.slowThreshold = threshold
	}
}

// logger is the main struct implementing the gormlogger.Interface.
type logger struct {
	manager                   *sklogger.Manager
	logLevel                  gormlogger.LogLevel
	slowThreshold             time.Duration
	ignoreRecordNotFoundError bool
}

// LogMode sets the log level for the logger and returns a new logger instance.
//
// Parameters:
//   - level: The gormlogger.LogLevel to set.
//
// Returns:
//   - A new gormlogger.Interface with the updated log level.
func (l *logger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.logLevel = level
	return &newLogger
}

// Info logs a message at Info level.
//
// Parameters:
//   - ctx: The context.Context for the log entry.
//   - msg: The message to log.
//   - data: Optional data to include in the log message.
func (l *logger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.logIfEnabled(ctx, gormlogger.Info, l.manager.Info, msg, data...)
}

// Warn logs a message at Warn level.
//
// Parameters:
//   - ctx: The context.Context for the log entry.
//   - msg: The message to log.
//   - data: Optional data to include in the log message.
func (l *logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.logIfEnabled(ctx, gormlogger.Warn, l.manager.Warn, msg, data...)
}

// Error logs a message at Error level.
//
// Parameters:
//   - ctx: The context.Context for the log entry.
//   - msg: The message to log.
//   - data: Optional data to include in the log message.
func (l *logger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.logIfEnabled(ctx, gormlogger.Error, l.manager.Error, msg, data...)
}

// logIfEnabled checks if logging is enabled for the given level and logs the message if it is.
//
// Parameters:
//   - ctx: The context.Context for the log entry.
//   - level: The gormlogger.LogLevel to check against.
//   - logFunc: The logging function to use (Info, Warn, or Error).
//   - msg: The message to log.
//   - data: Optional data to include in the log message.
func (l *logger) logIfEnabled(ctx context.Context, level gormlogger.LogLevel, logFunc func(context.Context, string, ...zap.Field), msg string, data ...interface{}) {
	if l.logLevel >= level {
		l.log(ctx, logFunc, msg, data...)
	}
}

// log is a generic logging function that formats the message and calls the appropriate log function.
//
// Parameters:
//   - ctx: The context.Context for the log entry.
//   - logFunc: The logging function to use (Info, Warn, or Error).
//   - msg: The message to log.
//   - data: Optional data to include in the log message.
func (l *logger) log(ctx context.Context, logFunc func(context.Context, string, ...zap.Field), msg string, data ...interface{}) {
	if len(data) > 0 {
		msg = fmt.Sprintf(msg, data...)
	}
	logFunc(ctx, msg)
}

// Trace logs the execution of SQL queries, including execution time and any errors.
//
// Parameters:
//   - ctx: The context.Context for the log entry.
//   - begin: The time when the SQL execution began.
//   - fc: A function that returns the SQL query and the number of rows affected.
//   - err: Any error that occurred during SQL execution.
func (l *logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.logLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()
	elapsedMs := fmt.Sprintf("%.3f ms", float64(elapsed.Nanoseconds())/1e6)

	// Determine the appropriate log level based on the execution result
	switch {
	case err != nil && l.logLevel >= gormlogger.Error && (!errors.Is(err, gormlogger.ErrRecordNotFound) || !l.ignoreRecordNotFoundError):
		// Log error if an error occurred and it's not an ignored "record not found" error
		l.manager.Error(ctx, "db error trace",
			zap.String("sql", sql),
			zap.Error(err),
			zap.String("elapsed", elapsedMs),
			zap.Int64("rows", rows),
		)
	case elapsed > l.slowThreshold && l.slowThreshold != 0 && l.logLevel >= gormlogger.Warn:
		// Log slow query warning if execution time exceeds the threshold
		l.manager.Warn(ctx, "db slow query",
			zap.String("sql", sql),
			zap.String("elapsed", elapsedMs),
			zap.Int64("rows", rows),
		)
	case l.logLevel >= gormlogger.Info:
		// Log general query information at Info level
		l.manager.Info(ctx, "db trace",
			zap.String("sql", sql),
			zap.String("elapsed", elapsedMs),
			zap.Int64("rows", rows),
		)
	}
}

// NewLog creates and returns a new logger instance with the given options.
//
// Parameters:
//   - manager: A pointer to the sklogger.Manager to use for logging.
//   - opts: A variadic list of LoggerOption functions to configure the logger.
//
// Returns:
//   - A new gormlogger.Interface configured with the provided options.
//
// Example:
//
//	logger := NewLog(manager, WithLevel("info"), WithSlowThreshold(300*time.Millisecond))
func NewLog(manager *sklogger.Manager, opts ...LoggerOption) gormlogger.Interface {
	l := &logger{
		manager:                   manager,
		logLevel:                  defaultLogLevel,
		slowThreshold:             defaultSlowThreshold,
		ignoreRecordNotFoundError: defaultIgnoreRecordNotFoundError,
	}

	// Apply all provided options
	for _, opt := range opts {
		opt(l)
	}

	return l
}
