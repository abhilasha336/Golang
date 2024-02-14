# Overview

This Go package provides logging functionalities to simplify logging operations within your application. It includes functions for setting log levels, including dump data, and making API requests to log external services.

## Index

- [InitLogger(clientOpt *ClientOptions, logType ...loggerImply) *Logger](#InitLogger)
- [Trace(message string, args ...interface{})](#Trace)
- [Debug(message string, args ...interface{})](#Debug)
- [Info(message string, args ...interface{})](#Info)
- [Warn(message string, args ...interface{})](#Warn)
- [Error(message string, args ...interface{})](#Error)
- [Fatal(message string, args ...interface{})](#Fatal)
- [Panic(message string, args ...interface{})](#Panic)
- [Tracef(message string, args ...interface{})](#Tracef)
- [Debugf(message string, args ...interface{})](#Debugf)
- [Infof(message string, args ...interface{})](#Infof)
- [Warnf(message string, args ...interface{})](#Warnf)
- [Errorf(message string, args ...interface{})](#Errorf)
- [Fatalf(message string, args ...interface{})](#Fatalf)
- [Panicf(message string, args ...interface{})](#Panicf)

## InitLogger

    InitLogger(clientOpt *ClientOptions, logType ...loggerImply) *Logger

This function is used to initialize the logger with a set of configurations.
oggerImply (Interface). `loggerImply` is an interface that defines a method for configuring client options for different `logging modes`. It is implemented by `FileMode` and `CloudMode` to customize the configuration based on the selected logging mode.

### Structures and Interfaces

### ClientOptions (Structure)

`ClientOptions` represents the overall client configuration for logging. It includes the following fields:
- `Service`: Name of the application or service generating the logs.
- `LogLevel`: Desired log level (e.g., "info", "debug", "error").
- `IncludeRequestDump`: Specifies whether to include request data in logs.
- `IncludeResponseDump`: Specifies whether to include response data in logs.
- `JsonFormater`: Enabling this option formats the output in a JSON structure. If set to 'false', the logs will default to a plain text format.

### FileMode (Structure)

`FileMode` represents the configuration for file-based 
logging mode. It includes the following fields:

- `LogPath`: Directory where log files will be stored.
- `LogfileName`: Format for generating log file names.
- `LogMaxSize`: Maximum size of a log file before rolling over.
- `LogMaxBackup`: Maximum number of old log files to retain.
- `LogMaxAge`: Maximum number of days to retain old log files.
- `LogMode`: Sets the log mode to save the logs in a file or console.
### CloudMode (Structure)

 `CloudMode` represents the configuration for cloud-based logging mode. It includes the following fields:
- `URL`: External service URL where logs will be sent.
- `Secret`: Client secret used to generate an authentication token.

## Logger Implementation Documentation

This documentation explains the implementation details of the logger service, including the structure, methods, and usage.

### Usage
This Go package offers easy-to-use logging functionalities for your application. Here is a brief overview of the main functions and how to use them:

Initialize logger by setting log level and configuring dump data.
To initialize the logger with a set of configurations, use the `InitLogger` function. Provide `ClientOptions` and additional logger options as needed.

Example: Initialize the logger

    // for logging to a file
    file := &logger.FileMode{
        LogfileName: "hello.log",
            .
            .
            .
        // add additional configuration
    }

    // for logging to database. For best practice pass URL and Secret via environment variables.
    db := &logger.CloudMode{
        URL:    "http://localhost:8000/api/v1.0/logs",
        Secret: "hello",
    }

    llog=logger.InitLogger(&logger.ClientOptions{
        Service:             consts.AppName,
        LogLevel:            "info",
        IncludeRequestDump:  true,
        IncludeResponseDump: true,
        JsonFormater:        false
    }, db, file)

   
    router := gin.Default() // or gin.New()

    //add the middleware in your service
    router.Use(middlewares.LogMiddleware(map[string]interface{}{}))

**Note: Pass request context in API calls.**


To log a basic error message, use the following code:

    Log().Error("This is an error message., err=%s", err.Error())
    //This logs an error message along with the error information retrieved from err.Error().

To include context information from a request context and log an error message, use:

    Log().WithContext(ctx.Request.Context()).Error("This is an error message. err=%s", err.Error())
    //This logs the error message along with the error information retrieved from err.Error(), providing additional context from the request context.

For more comprehensive logging with context and custom arguments, use the following code:

    Log().WithContext(ctx.Request.Context()).WithField(map[string]interface{}{"argument1":"value1", "argument2":"value2"}).Error("This is an error message.")
    //This logs the error message with the error information, alongside custom arguments provided as a map of key-value pairs.


# Logging Messages
The package provides functions for logging messages at different log levels:

## Trace 
Log a message at trace level.

## Debug
Log a message at debug level.

## Info
Log an informational message.

## Warn
Log a warning message.

## Error
Log an error message.

## Fatal
Log a message and exit with a non-zero status.

## Panic
Log a message and panic.

## Tracef
Log an error message with formatting at trace level

## Debugf
Log a message with formatting at debug level.

## Infof
Log an informational message with formatting.

## Warnf
Log a warning message with formatting. 

## Errorf
Log an error message with formatting.

## Fatalf
Log a message with formatting and exit with a non-zero status.

## Panicf
Log a message with formatting and panic.



## Sample Usage
```go
import (
    "gitlab.com/tuneverse/toolkit/core/logger"
)

// Initialize log file settings
file := &logger.FileMode{
    LogfileName:  "service_name.log", // Specify the log file name.
    LogPath:      "logs",            // Specify the log file path. Ensure this folder is added to .gitignore as "logs/".
    LogMaxAge:    7,                 // Set the maximum log file age in days.
    LogMaxSize:   1024 * 1024 * 10,  // Set the maximum log file size (10 MB).
    LogMaxBackup: 5,                 // Set the maximum number of log file backups to keep.
}

// Configure client options for the logger
clientOpt := &logger.ClientOptions{
    Service:             "demo", // Specify the service name.
    LogLevel:            "info", // Set the log level (e.g., "info", "error", "debug").
    IncludeRequestDump:  true,   // Include request data in logs.
    IncludeResponseDump: true,   // Include response data in logs.
    JsonFormater:        false   // log format, by default text format
}

// Initialize the logger
lg := logger.InitLogger(clientOpt, file)

// Usage examples with comments

// Log with context
lg.WithContext(ctx.Request.Context()).
    Error("Logging an error message with context")

// Log without context
lg.Error("Logging an error message without context")

// Alternatively, using the logger package directly

// Log with context
logger.Log().
    WithContext(ctx.Request.Context()).
    Error("Logging an error message with context")

// Log without context
logger.Log().
    Error("Logging an error message without context")

    // Log with context with formatting
logger.Log().
    WithContext(ctx.Request.Context()).
    Errorf("Logging an error message with context, err=%s", err.Error())

// Log without context with formatting
logger.Log().
    Errorf("Logging an error message without context, err=%s", err.Error())

```