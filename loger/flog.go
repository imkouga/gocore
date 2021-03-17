package loger

import (
	"errors"
	"fmt"
	"os"
	"time"
)

type fLogger struct {
	logger  *logger
	prefile *os.File
	file    *os.File
}

const (
	TRUNCATE_BY_DAY  int = 1
	TRUNCATE_BY_SIZE int = 2

	TRUNCATE_SIZE_128MB   = 1024 * 1024 * 128
	TRUNCATE_SIZE_256MB   = 1024 * 1024 * 256
	TRUNCATE_SIZE_512MB   = 1024 * 1024 * 512
	TRUNCATE_SIZE_1024MB  = 1024 * 1024 * 1024
	TRUNCATE_SIZE_DEFAULT = 1024 * 1024 * 512
)

var (
	flogger      *fLogger = nil
	fileName     string
	truncateMode int
	// default size is 512MB;
	defaultTruncateSize int64 = TRUNCATE_SIZE_DEFAULT
)

// 按照 年-月-日_时:分:秒 格式形式返回时间字符串
func GetCurrentTimeString() string {
	return time.Now().Format("2006-01-02_15:04:05")
}

func errorNewf(format string, v ...interface{}) error {
	text := fmt.Sprintf(format, v...)
	return errors.New(text)
}

func InitLogFile(file string, mode int, defaultSize int64) error {

	fName := fmt.Sprintf("%s-%s", file, GetCurrentTimeString())
	fd, err := os.OpenFile(fName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return errorNewf("Failed to openfile %s , error:[%s]", file, err.Error())
	}

	fileName = file
	truncateMode = mode
	if err = refreshSymlink(fd); err != nil {
		return err
	}

	l := newLoger(fd, "", Ldate|Ltime|Lshortfile, 1)
	fl := &fLogger{logger: l, prefile: nil, file: fd}
	flogger = fl

	switch truncateMode {
	case TRUNCATE_BY_DAY:
		go truncateByDay()
	case TRUNCATE_BY_SIZE:

		if defaultSize > 0 {
			defaultTruncateSize = defaultSize
		}

		go truncateBySize()
	default:
		return errors.New("invalid truncate mode.")
	}
	return nil
}

func isExist(filename string) bool {
	var err error
	if _, err = os.Open(filename); err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func refreshSymlink(file *os.File) error {
	if isExist(fileName) {
		os.Remove(fileName)
	}

	var err error
	if err = os.Symlink(file.Name(), fileName); err != nil {
		return errorNewf("symlink fail %s, error:[%s] ", fileName, err.Error())
	}

	return nil
}

func truncate() error {

	fName := fmt.Sprintf("%s-%s", fileName, GetCurrentTimeString())
	file, err := os.OpenFile(fName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return errorNewf("failed to open file. error[%s] ", err.Error())
	}

	flogger.logger.mu.Lock()
	defer flogger.logger.mu.Unlock()

	if flogger.prefile != nil {
		flogger.prefile.Close()
	}

	flogger.prefile = flogger.file
	flogger.file = file
	flogger.logger.out = file

	if err = refreshSymlink(flogger.file); err != nil {
		return err
	}

	return nil
}

func size() int64 {

	flogger.logger.mu.Lock()
	defer flogger.logger.mu.Unlock()

	if flogger.file == nil {
		return 0
	}

	fi, err := os.Stat(flogger.file.Name())
	if err != nil {
		return 0
	}

	return fi.Size()
}

func quit() {
	if flogger.prefile != nil {
		flogger.prefile.Close()
	}

	if flogger.file != nil {
		flogger.file.Close()
	}
}

func truncateBySize() {
	ticker := time.NewTicker(1 * time.Minute)
	defer func() {
		ticker.Stop()
		quit()
	}()

	for {
		select {
		case <-ticker.C:
			if size() >= defaultTruncateSize {
				if err := truncate(); err != nil {
					Errorf("failed to truncate, error[%s]", err.Error())
				}
			}
		}
	}
}

func truncateByDay() {
	ticker := time.NewTicker(1 * time.Minute)
	defer func() {
		ticker.Stop()
		quit()
	}()

	isNeedChange := true

	for {
		select {
		case <-ticker.C:
			now := time.Now()

			if now.Hour() == 0 && isNeedChange {

				isNeedChange = false
				if err := truncate(); err != nil {
					Errorf("failed to truncate. error[%s]", err.Error())
				}
				continue
			}

			if now.Hour() != 0 {
				isNeedChange = true
			}
		}
	}
}

func getFileLogger() *fLogger {
	return flogger
}

func isUninitFileLogger() bool {
	return getFileLogger() == nil
}

// Trace logs a message at trace level.
func FTrace(v ...interface{}) {

	if isUninitFileLogger() {
		Trace(v)
		return
	}

	if level <= LevelTrace {
		flogger.logger.Printf("[TRACE] %s", fmt.Sprintln(v...))
	}
}

// Debug logs a message at debug level.
func FDebug(v ...interface{}) {

	if isUninitFileLogger() {
		Debug(v)
		return
	}

	if level <= LevelDebug {
		flogger.logger.Printf("[DEBUG] %s", fmt.Sprintln(v...))
	}
}

// Info logs a message at info level.
func FInfo(v ...interface{}) {

	if isUninitFileLogger() {
		Info(v)
		return
	}

	if level <= LevelInfo {
		flogger.logger.Printf("[INFO] %s", fmt.Sprintln(v...))
	}
}

// Warning logs a message at warning level.
func FWarn(v ...interface{}) {

	if isUninitFileLogger() {
		Warn(v)
		return
	}

	if level <= LevelWarning {
		flogger.logger.Printf("[WARN] %s", fmt.Sprintln(v...))
	}
}

// Error logs a message at error level.
func FError(v ...interface{}) {

	if isUninitFileLogger() {
		Error(v)
		return
	}

	if level <= LevelError {
		flogger.logger.Printf("[ERROR] %s", fmt.Sprintln(v...))
	}
}

// Critical logs a message at critical level.
func FCritical(v ...interface{}) {

	if isUninitFileLogger() {
		Critical(v)
		return
	}

	if level <= LevelCritical {
		flogger.logger.Printf("[CRITICAL] %s", fmt.Sprintln(v...))
		os.Exit(1)
	}
}

// Tracef logs a message at trace level.
func FTracef(format string, v ...interface{}) {

	if isUninitFileLogger() {
		Tracef(format, v)
		return
	}

	message := fmt.Sprintf(format, v...)
	if level <= LevelTrace {
		flogger.logger.Printf("[TRACE] %s", fmt.Sprintln(message))
	}
}

// Debugf logs a message at debug level.
func FDebugf(format string, v ...interface{}) {

	if isUninitFileLogger() {
		Debugf(format, v)
		return
	}

	message := fmt.Sprintf(format, v...)
	if level <= LevelDebug {
		flogger.logger.Printf("[DEBUG] %s", fmt.Sprintln(message))
	}
}

// Infof logs a message at info level.
func FInfof(format string, v ...interface{}) {

	if isUninitFileLogger() {
		Infof(format, v)
		return
	}

	message := fmt.Sprintf(format, v...)
	if level <= LevelInfo {
		flogger.logger.Printf("[INFO] %s", fmt.Sprintln(message))
	}
}

// Warningf logs a message at warning level.
func FWarnf(format string, v ...interface{}) {

	if isUninitFileLogger() {
		Warnf(format, v)
		return
	}

	message := fmt.Sprintf(format, v...)
	if level <= LevelWarning {
		flogger.logger.Printf("[WARN] %s", fmt.Sprintln(message))
	}
}

// Errorf logs a message at error level.
func FErrorf(format string, v ...interface{}) {

	if isUninitFileLogger() {
		Errorf(format, v)
		return
	}

	message := fmt.Sprintf(format, v...)
	if level <= LevelError {
		flogger.logger.Printf("[ERROR] %s", fmt.Sprintln(message))
	}
}

// Criticalf logs a message at critical level.
func FCriticalf(format string, v ...interface{}) {

	if isUninitFileLogger() {
		Criticalf(format, v)
		return
	}

	message := fmt.Sprintf(format, v...)
	if level <= LevelCritical {
		flogger.logger.Printf("[CRITICAL] %s", fmt.Sprintln(message))
		os.Exit(1)
	}
}
