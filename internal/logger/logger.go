package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/sudhucodes/pingsystem/internal/config"
)

type Logger struct {
	mu          sync.Mutex
	logDir      string
	currentDate string
	file        *os.File
	logger      *log.Logger
}

var defaultLogger *Logger

// GetLogDir returns the log directory path (%APPDATA%\PingSystem\logs).
func GetLogDir() (string, error) {
	cfgDir, err := config.GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(cfgDir, "logs"), nil
}

// Init initializes the global logger.
func Init() error {
	logDir, err := GetLogDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	l := &Logger{
		logDir: logDir,
	}

	if err := l.rotateLogFile(); err != nil {
		return err
	}

	defaultLogger = l
	return nil
}

// rotateLogFile opens today's log file (YYYY-MM-DD.log) and creates multi-writer to stdout & logfile.
func (l *Logger) rotateLogFile() error {
	today := time.Now().Format("2006-01-02")
	if l.currentDate == today && l.file != nil {
		return nil
	}

	if l.file != nil {
		_ = l.file.Close()
	}

	fileName := fmt.Sprintf("%s.log", today)
	filePath := filepath.Join(l.logDir, fileName)

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file %s: %w", filePath, err)
	}

	l.file = file
	l.currentDate = today
	mw := io.MultiWriter(os.Stdout, file)
	l.logger = log.New(mw, "", log.LstdFlags)

	return nil
}

func (l *Logger) logMessage(level, format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	_ = l.rotateLogFile()

	msg := fmt.Sprintf(format, v...)
	fullMsg := fmt.Sprintf("[%s] %s", level, msg)

	if l.logger != nil {
		l.logger.Println(fullMsg)
	} else {
		log.Println(fullMsg)
	}
}

// Info logs an informational message.
func Info(format string, v ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.logMessage("INFO", format, v...)
	} else {
		log.Printf("[INFO] "+format+"\n", v...)
	}
}

// Warn logs a warning message.
func Warn(format string, v ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.logMessage("WARN", format, v...)
	} else {
		log.Printf("[WARN] "+format+"\n", v...)
	}
}

// Error logs an error message.
func Error(format string, v ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.logMessage("ERROR", format, v...)
	} else {
		log.Printf("[ERROR] "+format+"\n", v...)
	}
}

// Close closes the underlying log file.
func Close() {
	if defaultLogger != nil && defaultLogger.file != nil {
		_ = defaultLogger.file.Close()
	}
}
