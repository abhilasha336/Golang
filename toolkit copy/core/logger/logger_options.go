package logger

import (
	"log"

	"gitlab.com/tuneverse/toolkit/utils"
)

var (
	DefaultLogMaxAge      = 7
	DefaultSizeUnit       = 1024
	DefaultMinOne         = 1
	DefaultMinBackupCount = 3
	DefaultLogMaxSize     = int64(DefaultSizeUnit*DefaultSizeUnit) * 5
)

type logMode int

type loggerImply interface {
	setLogger(clientOpts *clientOptions) *clientOptions
}

type FileMode struct {
	// LogPath determines the directory in which to store log files.
	// It defaults to os.TempDir() if empty.
	LogPath string
	// LogfileName is the time formatting layout used to generate filenames.
	// It defaults to "2006-01-02T15-04-05.000.log".
	LogfileName string
	// LogMaxSize is the maximum size in bytes of the log file before it gets
	// rolled. It defaults to 5 megabytes.
	LogMaxSize int64
	// LogMaxBackup is the maximum number of old log files to retain. The default
	// is 3 files (though LogMaxAge may still cause them to get deleted.)
	LogMaxBackup int
	// LogMaxAge is the maximum number of days to retain old log files based on
	// FileInfo.ModTime. Note that a day is defined as 24 hours and may not
	// exactly correspond to calendar days due to daylight savings, leap seconds, etc.
	LogMaxAge int
}

type CloudMode struct {
	// URL of the external service to send the logs
	URL string
	// Secret is the client secret to generate the token
	Secret string
}

// set the client options for file mode
func (file *FileMode) setLogger(clientOpts *clientOptions) *clientOptions {
	if utils.IsEmpty(file.LogPath) {
		clientOpts.logPath = utils.TempDir()
	} else {
		clientOpts.logPath = file.LogPath
	}
	if !utils.IsEmpty(file.LogfileName) {
		clientOpts.logfileName = file.LogfileName
	}
	if file.LogMaxSize < DefaultLogMaxSize {
		clientOpts.logMaxSize = DefaultLogMaxSize
	} else {
		clientOpts.logMaxSize = file.LogMaxSize
	}
	if file.LogMaxAge < DefaultMinOne {
		clientOpts.logMaxAge = DefaultLogMaxAge
	} else {
		clientOpts.logMaxAge = file.LogMaxAge
	}
	if file.LogMaxBackup < DefaultMinOne {
		clientOpts.logMaxBackup = DefaultMinBackupCount
	} else {
		clientOpts.logMaxBackup = file.LogMaxBackup
	}
	clientOpts.appendLogModes(logModeFile)
	return clientOpts
}

// set the client options for cloud mode
func (cloud *CloudMode) setLogger(clientOpts *clientOptions) *clientOptions {
	var err error
	clientOpts.recordOptions, err = initTransportOptions(true, cloud.URL, cloud.Secret)
	if err != nil {
		log.Fatalf("SetTransportOptions failed error=%s", err.Error())
	}
	clientOpts.appendLogModes(logModeCloud)
	return clientOpts
}
