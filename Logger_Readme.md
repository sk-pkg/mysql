# MySQL Logger Package Usage Documentation

## Introduction

This package provides a custom MySQL logger that implements GORM's `logger.Interface`. This logger integrates with the `sk-pkg/logger` package, offering flexible configuration options and detailed SQL query logging.

## Installation

Ensure you have the necessary dependencies installed in your project:

```shell
go get github.com/sk-pkg/logger
go get gorm.io/gorm
```

## Usage

### Creating a New Logger

Use the `NewLog` function to create a new logger instance:

```go
import (
    sklogger "github.com/sk-pkg/logger"
    "gorm.io/gorm"
    mysqlLogger "your/package/path"
)

// Assuming you already have a sklogger.Manager instance
manager := // ... initialize sklogger.Manager

// Create a new logger
logger := mysqlLogger.NewLog(manager)

// Use this logger in GORM configuration
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
    Logger: logger,
})
```

### Configuration Options

You can customize the logger's behavior using the following options:

#### 1. Setting Log Level

Use the `WithLevel` option to set the log level:

```go
logger := mysqlLogger.NewLog(manager, mysqlLogger.WithLevel("info"))
```

Available log levels are:
- "debug"
- "info"
- "warn"
- "silent"

#### 2. Ignoring "Record Not Found" Errors

Use the `WithIgnoreRecordNotFoundError` option to set whether to ignore "record not found" errors:

```go
logger := mysqlLogger.NewLog(manager, mysqlLogger.WithIgnoreRecordNotFoundError(true))
```

#### 3. Setting Slow Query Threshold

Use the `WithSlowThreshold` option to set the time threshold for slow queries:

```go
logger := mysqlLogger.NewLog(manager, mysqlLogger.WithSlowThreshold(500 * time.Millisecond))
```

### Combining Multiple Options

You can combine multiple options to configure the logger:

```go
logger := mysqlLogger.NewLog(manager,
    mysqlLogger.WithLevel("info"),
    mysqlLogger.WithIgnoreRecordNotFoundError(true),
    mysqlLogger.WithSlowThreshold(300 * time.Millisecond),
)
```

## Log Output

This logger will output different types of logs based on the configured log level and query execution time:

1. Error logs: When a query fails (unless it's an ignored "record not found" error).
2. Slow query warnings: When query execution time exceeds the set threshold.
3. General query information: Logs basic information for all queries (if log level is set to Info or lower).

Each log entry includes the following information:
- SQL query statement
- Execution time
- Number of rows affected
- Error message (if any)

## Considerations

- Ensure appropriate log levels are set in production environments to avoid excessive log output that could impact performance.
- The slow query threshold should be set based on your application requirements and database performance.
- Using `WithIgnoreRecordNotFoundError(true)` can reduce unnecessary "record not found" error logs, but ensure this doesn't affect your error handling logic.