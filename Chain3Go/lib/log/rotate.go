package log

import (
	"fmt"
	"os"
	"reflect"
	"time"
)

// LogRotate: change file handler by date
func LogRotate(prefixFile string, lvl Lvl) {
	logger := Root()
	glogger := NewGlogHandler(StreamHandler(os.Stderr, TerminalFormat(false)))
	glogger.Verbosity(lvl)
	lastFileHandler, _ := setLogFileHandler(prefixFile, lvl, glogger, logger, true)
	for {
		now := time.Now()
		// change log file everyday
		next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		t := time.NewTimer(next.Sub(now))
		<-t.C
		// change log file to date
		if temp, err := setLogFileHandler(prefixFile, lvl, glogger, logger, false); err == nil {
			reflect.ValueOf(lastFileHandler).Field(0).MethodByName("Close").Call([]reflect.Value{})
			lastFileHandler = temp
		}
	}
}

// setLogFileHandler: and return file handler.
func setLogFileHandler(prefixFile string, lvl Lvl, glogger *GlogHandler, logger Logger, bMust bool) (Handler, error) {
	curFilename := getLogFileName(prefixFile, time.Now())
	var fileHandler Handler
	var err error
	if bMust {
		fileHandler = Must.FileHandler(curFilename, TerminalFormat(false))
	} else {
		fileHandler, err = FileHandler(curFilename, TerminalFormat(false))
	}

	if err == nil {
		logger.SetHandler(MultiHandler(
			LvlFilterHandler(lvl, fileHandler),
			glogger))
	}
	return fileHandler, err
}

func getLogFileName(prefixFile string, t time.Time) string {
	return fmt.Sprintf("%v-%04d-%02d-%02d.log", prefixFile, t.Year(), t.Month(), t.Day())
}
