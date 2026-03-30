package logger

import (
	"fmt"
	"io"
	stdlog "log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type AppLogger struct {
	*stdlog.Logger
	errorLogFile *os.File
}

// New creates a new AppLogger instance.
// It initializes the logger to write to both stderr and a log file in the specified log directory.
func New(logDir string) (*AppLogger, error) {
	// Open error log file
	errorLogFile, err := openErrorLogFile(logDir)
	if err != nil {
		return nil, err
	}

	// Create a multi-writer to write to both stderr and the error log file
	errorWriter := io.MultiWriter(os.Stderr, errorLogFile)

	// Set Gin's default error writer to our multi-writer so that Gin's internal logs also go to the same destinations
	gin.DefaultErrorWriter = errorWriter

	return &AppLogger{
		Logger:       stdlog.New(errorWriter, "[app] ", stdlog.LstdFlags|stdlog.Lshortfile),
		errorLogFile: errorLogFile,
	}, nil
}

// Close closes the error log file if it is open.
func (l *AppLogger) Close() error {
	if l == nil || l.errorLogFile == nil {
		return nil
	}

	return l.errorLogFile.Close()
}

// CloseWithReport closes the error log file and reports any error to stderr.
func (l *AppLogger) CloseWithReport() {
	if err := l.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "close error log file: %v\n", err)
	}
}

// openErrorLogFile opens the error log file in the specified log directory.
// It creates the log directory if it does not exist and returns the opened file.
func openErrorLogFile(logDir string) (*os.File, error) {
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return nil, fmt.Errorf("create log dir %q: %w", logDir, err)
	}

	logPath := filepath.Join(logDir, "error.log")
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("open error log file %q: %w", logPath, err)
	}

	return file, nil
}
