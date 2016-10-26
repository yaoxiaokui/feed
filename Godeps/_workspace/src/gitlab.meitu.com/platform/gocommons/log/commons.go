// Package log define common log interface, also with a default implementation which wrap seelog.
package log

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"sync"
)

type loggerCreatorType func(name string) (Logger, error)

var loggerCreator loggerCreatorType = func(name string) (Logger, error) {
	return EmptyLogger, nil
}

var defaultStackDepth = 2

var loggerCache map[string]Logger = make(map[string]Logger)

var m *sync.Mutex = new(sync.Mutex)

// GetLoggerByConfigFile return Logger instance init by specified config file
func GetLoggerByConfigFile(path string) Logger {
	if !filepath.IsAbs(path) {
		var err error
		path, err = filepath.Abs(path)
		if err != nil {
			fmt.Printf("can not resolve logger config file %s: %s\n", path, err.Error())
			return EmptyLogger
		}
	}
	return GetLogger(path)
}

// GetLogger return Logger instance init by $application/conf/seelog-$name.xml config file
func GetLogger(name string) Logger {
	logger, exists := loggerCache[name]

	if exists == false {
		m.Lock()
		defer m.Unlock()

		logger, exists = loggerCache[name]
		if exists == false {
			var err error
			logger, err = loggerCreator(name)
			if err != nil {
				fmt.Printf("create logger %v failed. %s\n", name, err.Error())
				logger = EmptyLogger
			}
			logger.SetAdditionalStackDepth(defaultStackDepth)
			loggerCache[name] = logger
		}
	}

	return logger
}

func FlushAll() {
	for _, logger := range loggerCache {
		logger.Flush()
	}
}

var EmptyLogger = &emptyLogger{}

//空的logger对象
type emptyLogger struct {
}

func (logger *emptyLogger) Debug(v ...interface{}) {

}
func (logger *emptyLogger) Info(v ...interface{}) {

}
func (logger *emptyLogger) Error(v ...interface{}) {

}
func (logger *emptyLogger) Warn(v ...interface{}) {

}
func (logger *emptyLogger) Fatal(v ...interface{}) {

}

func (logger *emptyLogger) Debugf(format string, params ...interface{}) {

}
func (logger *emptyLogger) Infof(format string, params ...interface{}) {

}
func (logger *emptyLogger) Errorf(format string, params ...interface{}) {

}
func (logger *emptyLogger) Warnf(format string, params ...interface{}) {

}
func (logger *emptyLogger) Fatalf(format string, params ...interface{}) {
}
func (logger *emptyLogger) SetAdditionalStackDepth(depth int) {
}
func (logger *emptyLogger) Flush() {

}

func getCallerInfo(skip int) (string, int) {
	if skip == 0 {
		return "???", 0
	}

	var (
		file string
		line int
		ok   bool
	)

	for {
		_, file, line, ok = runtime.Caller(skip + 1)
		if !ok {
			return "???", 0
		}

		if file != "<autogenerated>" {
			break
		}

		skip += 1
	}

	return string(file[len(filepath.Dir(file))+1:]), line
}

// current is the logger used in all package level convenience funcs like 'Debug', 'Info', 'Error', 'Flush', etc.
var current Logger

func init() {
	current = &consoleLogger{}
	current.SetAdditionalStackDepth(3)
}

func UseLoggerByName(name string) error {
	logger := GetLogger(name)
	if logger == nil {
		return fmt.Errorf("can not get logger %v", name)
	}
	if current != nil {
		current.Flush()
	}
	current = logger

	return nil
}

func UseLogger(logger Logger) error {
	if logger == nil {
		return errors.New("logger can not be nil")
	}

	if current != nil {
		current.Flush()
	}
	current = logger

	return nil
}

func SetAdditionalStackDepth(depth int) {
	current.SetAdditionalStackDepth(depth)
}

func Debug(v ...interface{}) {
	current.Debug(v...)
}
func Info(v ...interface{}) {
	current.Info(v...)
}
func Error(v ...interface{}) {
	current.Error(v...)
}
func Warn(v ...interface{}) {
	current.Warn(v...)
}
func Fatal(v ...interface{}) {
	current.Fatal(v...)
}

func Debugf(format string, args ...interface{}) {
	current.Debugf(format, args...)
}
func Infof(format string, args ...interface{}) {
	current.Infof(format, args...)
}
func Errorf(format string, args ...interface{}) {
	current.Errorf(format, args...)
}
func Warnf(format string, args ...interface{}) {
	current.Warnf(format, args...)
}
func Fatalf(format string, args ...interface{}) {
	current.Fatalf(format, args...)
}

func Flush() {
	current.Flush()
}
