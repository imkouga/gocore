package loger

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type rLogger struct {
	logger *logger
	file   *os.File
}

var (
	rlogger     *rLogger = nil
	logrotateFn string
)

func InitLogrotateFile(filename string) error {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return errors.New("failed to open file " + filename + ". - " + err.Error())
	}

	l := newLoger(file, "", Ldate|Ltime|Lshortfile, 1)
	fl := &rLogger{logger: l, file: file}
	rlogger = fl
	logrotateFn = filename

	go logrotate()
	return nil
}

func handleSignal() {
	rlogger.logger.mu.Lock()
	defer rlogger.logger.mu.Unlock()

	if rlogger.file != nil {
		rlogger.file.Close()
	}

	file, err := os.OpenFile(logrotateFn, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		var err1 error
		file, err1 = os.OpenFile(logrotateFn, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		if err1 != nil {
			Error("failed to open|create file " + logrotateFn + ". - " + err1.Error())
			return
		}
	}

	rlogger.file = file
	rlogger.logger.out = file

	Info("success to truncate log file.")
}

func logrotate() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGUSR1)

	for {
		select {
		case sig := <-sigs:
			Info("get signal:", sig)
			handleSignal()
		}
	}
}

// Trace logs a message at trace level.
func RTrace(v ...interface{}) {
	if level <= LevelTrace {
		rlogger.logger.Printf("[TRACE] %s", fmt.Sprintln(v...))
	}
}

// Debug logs a message at debug level.
func RDebug(v ...interface{}) {
	if level <= LevelDebug {
		rlogger.logger.Printf("[DEBUG] %s", fmt.Sprintln(v...))
	}
}

// Info logs a message at info level.
func RInfo(v ...interface{}) {
	if level <= LevelInfo {
		rlogger.logger.Printf("[INFO] %s", fmt.Sprintln(v...))
	}
}

// Warning logs a message at warning level.
func RWarn(v ...interface{}) {
	if level <= LevelWarning {
		rlogger.logger.Printf("[WARN] %s", fmt.Sprintln(v...))
	}
}

// Error logs a message at error level.
func RError(v ...interface{}) {
	if level <= LevelError {
		rlogger.logger.Printf("[ERROR] %s", fmt.Sprintln(v...))
	}
}

// Critical logs a message at critical level.
func RCritical(v ...interface{}) {
	if level <= LevelCritical {
		rlogger.logger.Printf("[CRITICAL] %s", fmt.Sprintln(v...))
		os.Exit(1)
	}
}
