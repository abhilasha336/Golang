package utilities

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

// NewLogger function which logs error with function name and errors
func NewLogger(fName string) *log.Entry {
	return log.WithFields(log.Fields{
		"fn": fmt.Sprintf("%s()", fName),
	})
}
