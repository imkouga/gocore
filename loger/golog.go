package loger

import (
	"fmt"
	"io"
	"os"
	"strings"
)

//--------------------
// LOG LEVEL
//--------------------

// Log levels to control the logging output.
const (
	LevelTrace = iota
	LevelDebug
	LevelInfo
	LevelWarning
	LevelError
	LevelCritical
)

// logLevel controls the global log level used by the logger.
var (
	level        = LevelTrace
	logPrefixStr string
)

// LogLevel returns the global log level and can be used in
// own implementations of the logger interface.
func GetLevel() int {
	return level
}

// SetLogLevel sets the global log level used by the simple
// logger.
func SetLevel(l int) {
	level = l
}

func SetLogLinePreFix(prefix string) {
	prefix = fmt.Sprintf("%s ", strings.TrimSpace(prefix))
	goLogger.prefix = prefix
}

// SetFlags sets the output flags for the logger.
func SetFlags(flag int) {
	goLogger.SetFlags(flag)
}

// logger references the used application logger.
var goLogger = newLoger(os.Stdout, "", Ldate|Ltime|Lshortfile, 1)

func SetOutput(w io.Writer) {
	goLogger = newLoger(w, "", Ldate|Ltime|Lshortfile, 1)
}

func SetLogFullFileLine(clearZone string) {
	goLogger.flag = Ldate | Ltime | Llongfile
	goLogger.clearPrefix = clearZone
}

// SetLogger sets a new logger.
func SetLogger(l *logger) {
	goLogger = l
}

// Trace logs a message at trace level.
func Trace(v ...interface{}) {
	if level <= LevelTrace {
		goLogger.Printf("[TRACE] %s %s", logPrefixStr, fmt.Sprintln(v...))
	}
}

// Debug logs a message at debug level.
func Debug(v ...interface{}) {
	if level <= LevelDebug {
		goLogger.Printf("[DEBUG] %s %s", logPrefixStr, fmt.Sprintln(v...))
	}
}

// Info logs a message at info level.
func Info(v ...interface{}) {
	if level <= LevelInfo {
		goLogger.Printf("[INFO] %s %s", logPrefixStr, fmt.Sprintln(v...))
	}
}

// Warning logs a message at warning level.
func Warn(v ...interface{}) {
	if level <= LevelWarning {
		goLogger.Printf("[WARN] %s %s", logPrefixStr, fmt.Sprintln(v...))
	}
}

// Error logs a message at error level.
func Error(v ...interface{}) {
	if level <= LevelError {
		goLogger.Printf("[ERROR] %s %s", logPrefixStr, fmt.Sprintln(v...))
	}
}

// Critical logs a message at critical level.
func Critical(v ...interface{}) {
	if level <= LevelCritical {
		goLogger.Printf("[CRITICAL] %s %s", logPrefixStr, fmt.Sprintln(v...))
		os.Exit(1)
	}
}

// ----------------------------------------------------------------------------------------------

// Trace logs a message at trace level.
func Tracef(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	if level <= LevelTrace {
		goLogger.Printf("[TRACE] %s %s", logPrefixStr, fmt.Sprintln(message))
	}
}

// Debug logs a message at debug level.
func Debugf(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	if level <= LevelDebug {
		goLogger.Printf("[DEBUG] %s %s", logPrefixStr, fmt.Sprintln(message))
	}
}

// Info logs a message at info level.
func Infof(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	if level <= LevelInfo {
		goLogger.Printf("[INFO] %s %s", logPrefixStr, fmt.Sprintln(message))
	}
}

// Warning logs a message at warning level.
func Warnf(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	if level <= LevelWarning {
		goLogger.Printf("[WARN] %s %s", logPrefixStr, fmt.Sprintln(message))
	}
}

// Error logs a message at error level.
func Errorf(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	if level <= LevelError {
		goLogger.Printf("[ERROR] %s %s", logPrefixStr, fmt.Sprintln(message))
	}
}

// Critical logs a message at critical level.
func Criticalf(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	if level <= LevelCritical {
		goLogger.Printf("[CRITICAL] %s %s", logPrefixStr, fmt.Sprintln(message))
		os.Exit(1)
	}
}
