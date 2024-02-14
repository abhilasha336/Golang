package logger

import (
	"slices"

	"github.com/sirupsen/logrus"
)

type (
	clientOptions struct {
		service      string
		logPath      string
		logfileName  string
		logMaxSize   int64
		logMaxBackup int
		logMaxAge    int
		logLevel     logrus.Level

		//include request dump data
		includeRequestDump bool

		//include response dump
		includeResponseDump bool

		//record the logs in external service
		recordOptions *recordOptions

		//set the log mode to save the logs in a file or console
		//by default it will set to save logs in console
		logMode []logMode
	}

	recordOptions struct {
		url    string
		secret string
		token  string
	}
)

func (clientOpt *clientOptions) appendLogModes(mode logMode) {
	if !slices.Contains(clientOpt.logMode, mode) {
		clientOpt.logMode = append(clientOpt.logMode, mode)
	}
}

// set request dump in log request
func (clientOpt *clientOptions) setRequestData(includeDump bool) {
	clientOpt.includeRequestDump = includeDump
}

// set response dump in log request
func (clientOpt *clientOptions) setResponseData(includeResponseDump bool) {
	clientOpt.includeResponseDump = includeResponseDump
}

// log level - possible values trace, debug, info, warn, error, fatal, panic
func (clientOpt *clientOptions) setLogLevel(logLevel string) {
	clientOpt.logLevel, _ = logrus.ParseLevel(logLevel)
}
func (clientOpt *clientOptions) setRequestDumpData() {
	includeRequestDumpData = clientOpt.includeRequestDump
}

func (clientOpt *clientOptions) setServiceName() {
	service = clientOpt.service
}
func (clientOpt *clientOptions) setResponseDumpData() {
	includeResponseDumpData = clientOpt.includeResponseDump
}
