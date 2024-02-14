package logger

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"gitlab.com/tuneverse/toolkit/consts"
	"gitlab.com/tuneverse/toolkit/utils"
	"gopkg.in/natefinch/lumberjack.v1"
)

var (
	recordLog               bool
	logToken                string
	logServiceURL           string
	logger                  *logrus.Logger
	includeRequestDumpData  bool
	includeResponseDumpData bool
	service                 string
)

const (
	logModeFile logMode = iota
	logModeCloud
)

type ClientOptions struct {
	// Service describes the application name to be logged
	Service string

	//possible values are trace, debug, info, warn, error, fatal, panic
	//defaults panic level
	LogLevel string

	//include request dump data
	IncludeRequestDump bool

	//include response dump
	IncludeResponseDump bool

	//JsonFormater
	JSONFormater bool
}

type Logger struct {
	ctx        context.Context
	entry      *logrus.Entry
	fields     map[string]interface{}
	jsonFormat bool
	mu         *sync.RWMutex
}

var logObject *Logger

func Log() *Logger {
	return logObject
}

// logger initialization
func InitLogger(clientOpt *ClientOptions, logType ...loggerImply) *Logger {
	clientOpts := &clientOptions{}
	if utils.IsEmpty(clientOpt.Service) {
		log.Fatal("service name is required")
	}

	for _, logger := range logType {
		clientOpts = logger.setLogger(clientOpts)
	}
	clientOpts.service = clientOpt.Service
	clientOpts.setRequestData(clientOpt.IncludeRequestDump)
	clientOpts.setResponseData(clientOpt.IncludeResponseDump)
	clientOpts.setLogLevel(clientOpt.LogLevel)

	logger = logrus.New()
	clientOpts.setRequestDumpData()
	clientOpts.setResponseDumpData()
	clientOpts.setServiceName()

	for _, v := range clientOpts.logMode {
		err := v.init(clientOpts)
		if err != nil {
			log.Fatal("logger initialisation failed %w", err)
		}
	}

	logObject = &Logger{
		fields:     make(map[string]interface{}),
		jsonFormat: clientOpt.JSONFormater,
		mu:         &sync.RWMutex{},
	}
	return logObject
}

func (lm logMode) init(cl *clientOptions) error {

	logger.SetLevel(cl.logLevel)
	logger.SetFormatter(&logrus.TextFormatter{})

	switch lm {
	case logModeCloud:
		cl.recordOptions.setLogServiceCredentials()
	case logModeFile:
		if _, err := os.Stat(cl.logPath); os.IsNotExist(err) {
			err := os.MkdirAll(cl.logPath, os.ModePerm)
			if err != nil {
				return err
			}
		}
		logger.SetOutput(io.MultiWriter(&lumberjack.Logger{
			NameFormat: cl.logfileName,
			Dir:        cl.logPath,
			MaxSize:    cl.logMaxSize,
			MaxBackups: cl.logMaxBackup,
			MaxAge:     cl.logMaxAge,
		}, os.Stdout))
	}
	return nil
}

func initTransportOptions(recordLogs bool, url, tokenSecret string) (*recordOptions, error) {
	transport := &recordOptions{}
	if recordLogs {
		if utils.IsEmpty(url) || utils.IsEmpty(tokenSecret) {
			return transport, fmt.Errorf("invalid transport options")
		}
		//generating token without expiration time to reuse the token for every request
		token, err := utils.GenerateJWTAuthToken(tokenSecret, map[string]interface{}{})
		if err != nil {
			return transport, fmt.Errorf("token generation failed, err=%s", err.Error())
		}
		//initial ping to the logger service
		err = ping(fmt.Sprintf("%s/%s", url, "health"), token)
		if err != nil {
			return transport, err
		}
		transport.url = url
		transport.secret = tokenSecret
		transport.token = token

	}
	return transport, nil

}

// ping request
func ping(url, token string) error {

	headers := map[string]interface{}{
		"Authorization": token,
	}
	resp, err := utils.APIRequest(http.MethodGet, url, headers, nil)
	if err != nil {
		return fmt.Errorf("ping request failed: %s", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ping request failed: invalid transport Options")
	}
	return nil
}

func (rc *recordOptions) setLogServiceCredentials() {
	recordLog = true
	logServiceURL = rc.url
	logToken = rc.token
}

func GetRequestDumpStatus() bool {
	return includeRequestDumpData
}

func GetResponseDumpStatus() bool {
	return includeResponseDumpData
}

func GetService() string {
	return service
}

// Get the context
func (log *Logger) Context() context.Context {
	return log.ctx
}

// With Context
func (log Logger) WithContext(ctx context.Context) Logger {
	log.ctx = ctx
	return log
}

// Entry
func (log Logger) Entry() *logrus.Entry {
	return log.entry
}

// WithEntry
func (log *Logger) WithEntry(entry *logrus.Entry) *Logger {
	log.entry = entry
	return log
}

// Fields
func (log *Logger) Fields() map[string]interface{} {
	return log.fields
}

// logFunc
func (log Logger) logFunc(fFunc bool, fn string, message string, args ...interface{}) {
	log.mu.Lock()
	defer log.mu.Unlock()

	// Get the level from the fn name
	level, err := logrus.ParseLevel(fn)
	if err != nil {
		level = logrus.FatalLevel
	}

	// get the existing fields
	var fields = log.Fields()

	fields[consts.ContextMessage] = message

	if fFunc {
		fields[consts.ContextMessage] = fmt.Sprintf(message, args...)
	} else {
		fields["args"] = fmt.Sprint(args...)
	}

	obje := log.WithFields(fields).Log(level).Entry()

	callableMethod := reflect.ValueOf(obje).MethodByName(fn)
	if callableMethod.IsValid() {
		inputs := make([]reflect.Value, 0)
		inputs = append(inputs, reflect.ValueOf(message))

		for _, v := range args {
			inputs = append(inputs, reflect.ValueOf(v))
		}

		callableMethod.Call(inputs)
		return
	}
}

// captureLogs
func (log Logger) captureLogs(ctx context.Context, fields map[string]interface{}, incomingLogLevel logrus.Level) {
	if logger.Level >= incomingLogLevel && recordLog {
		logg := logrus.New()
		headers := map[string]interface{}{
			"Authorization": logToken,
		}
		resp, err := utils.APIRequest(http.MethodPost, fmt.Sprintf("%s/%s", logServiceURL, "logs"), headers, fields)
		if err != nil {
			logg.Errorf("capturing logs failed, api call failed, err=%s", err.Error())
		}
		if resp.StatusCode != http.StatusCreated {
			logg.Errorf("capturing logs failed, api responsecode=%v", resp.StatusCode)
		}
	}
}

// WithFields
func (log Logger) WithFields(fields map[string]interface{}) Logger {
	log.fields = fields
	return log
}

// Log
func (log Logger) Log(incomingLogLevel logrus.Level) Logger {
	fields := log.Fields()
	ctx := log.Context()

	if ctx != nil {
		if ctxFields, ok := ctx.Value(consts.LogData).(map[string]interface{}); ok {
			for key, value := range ctxFields {
				fields[key] = value
			}
		}
	}

	lf := make(logrus.Fields)
	entry := logger.WithFields(lf)
	for key, element := range fields {
		entry = entry.WithField(key, element)
	}
	fields[consts.ContextTimeStamp] = time.Now()
	fields[consts.ContextLogLevel] = incomingLogLevel.String()
	log.captureLogs(ctx, fields, incomingLogLevel)

	// enable to set it as json formater
	if log.jsonFormat {
		entry.Logger.SetFormatter(&logrus.JSONFormatter{})
	}

	// log entry
	log.WithEntry(entry)

	return log
}

// Trace
func (log Logger) Trace(message string, args ...interface{}) {
	log.logFunc(false, "Trace", message, args...)
}

// Tracef
func (log Logger) Tracef(message string, args ...interface{}) {
	log.logFunc(true, "Trace", message, args...)
}

// Errorf
func (log Logger) Errorf(message string, args ...interface{}) {
	var fields = log.Fields()
	_, file, no, ok := runtime.Caller(1)
	if ok {
		fName := strings.Split(file, "/")
		fields["file"] = fName[len(fName)-1]
		fields["line"] = no
		fields[consts.ContextMessage] = fmt.Sprintf(message, args...)
	}

	// set the message
	fields[consts.ContextMessage] = fmt.Sprintf(message, args...)

	log.
		WithFields(fields).
		Log(logrus.ErrorLevel).
		Entry().Errorf(message, args...)

}

// Error
func (log Logger) Error(message string, args ...interface{}) {

	var fields = log.Fields()
	_, file, no, ok := runtime.Caller(1)
	if ok {
		fName := strings.Split(file, "/")
		fields["file"] = fName[len(fName)-1]
		fields["line"] = no
		fields[consts.ContextMessage] = message
		fields["args"] = fmt.Sprint(args...)
	}

	fields[consts.ContextMessage] = message

	log.
		WithFields(fields).
		Log(logrus.ErrorLevel).
		Entry().Error(args...)
}

// Print
func (log Logger) Print(message string, args ...interface{}) {
	log.logFunc(false, "Print", message, args...)
}

// Printf
func (log Logger) Printf(message string, args ...interface{}) {
	log.logFunc(true, "Print", message, args...)
}

// Info
func (log Logger) Info(message string, args ...interface{}) {
	log.logFunc(false, "Info", message, args...)
}

// Infof
func (log Logger) Infof(message string, args ...interface{}) {
	log.logFunc(true, "Info", message, args...)
}

// Debug
func (log Logger) Debug(message string, args ...interface{}) {
	log.logFunc(false, "Debug", message, args...)
}

// Debugf
func (log Logger) Debugf(message string, args ...interface{}) {
	log.logFunc(true, "Debug", message, args...)
}

// Warn
func (log Logger) Warn(message string, args ...interface{}) {
	log.logFunc(false, "Warn", message, args...)
}

// Warnf
func (log Logger) Warnf(message string, args ...interface{}) {
	log.logFunc(true, "Warn", message, args...)
}

// Fatal
func (log Logger) Fatal(message string, args ...interface{}) {
	log.logFunc(false, "Fatal", message, args...)
}

// Fatalf
func (log Logger) Fatalf(message string, args ...interface{}) {
	log.logFunc(true, "Fatalf", message, args...)
}

// Panic
func (log Logger) Panic(message string, args ...interface{}) {
	log.logFunc(false, "Panic", message, args...)
}

// Panicf
func (log Logger) Panicf(message string, args ...interface{}) {
	log.logFunc(true, "Panicf", message, args...)
}
