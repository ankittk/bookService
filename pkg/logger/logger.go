package logger

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"sort"
	"sync"
	"time"
)

var defaultLogger *Logger

func init() {
	defaultLogger = NewLogger()
}

type Logger struct {
	logger *log.Logger
	tags   []slog.Attr
	mu     sync.Mutex // Mutex to guard the tags slice
}

type LogLevel int

const (
	InfoLevel LogLevel = iota
	DebugLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

func NewLogger() *Logger {
	return &Logger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
		tags:   []slog.Attr{},
	}
}

func Info(ctx context.Context, message any) *Logger {
	return defaultLogger.Info(ctx, message)
}

func Warn(ctx context.Context, message any) *Logger {
	return defaultLogger.Warn(ctx, message)
}

func Error(ctx context.Context, message any) *Logger {
	return defaultLogger.Error(ctx, message)
}

func Debug(ctx context.Context, message any) *Logger {
	return defaultLogger.Debug(ctx, message)
}

func Fatal(ctx context.Context, message any) *Logger {
	return defaultLogger.Fatal(ctx, message)
}

func Add(attrs ...slog.Attr) *Logger {
	return defaultLogger.add(attrs...)
}

// Custom function to extract context values dynamically from a map
func (l *Logger) extractContextValues(ctx context.Context) string {
	var contextDetails string

	// Retrieve the value from the context (assuming it's a map)
	if data, ok := ctx.Value("context_data").(map[string]interface{}); ok {
		var keys []string
		for key := range data {
			keys = append(keys, key)
		}

		sort.Strings(keys) // Sort the keys lexicographically

		// Print context data in sorted order
		for _, key := range keys {
			contextDetails += fmt.Sprintf("[%s: %v] ", key, data[key])
		}
	}

	return contextDetails
}

func (l *Logger) logInternal(ctx context.Context, level LogLevel, message any) {
	contextDetails := l.extractContextValues(ctx)

	// Format the log message with the current timestamp
	logMessage := fmt.Sprintf("[%s] %s%s", time.Now().Format(time.RFC3339), contextDetails, message)

	// Add tags to the log message if available
	if len(l.tags) > 0 {
		logMessage += " | Tags: "
		// Iterate over tags
		for _, attr := range l.tags {
			logMessage += fmt.Sprintf("%s=%v ", attr.Key, attr.Value)
		}
	}

	// Print with color depending on the log level
	color := getColorForLogLevel(level)
	logMessage = fmt.Sprintf("%s%s\033[0m", color, logMessage)

	// Output the message
	if level >= ErrorLevel {
		l.logger.SetFlags(0)
		l.logger.SetOutput(os.Stderr)
	} else {
		l.logger.SetFlags(0)
		l.logger.SetOutput(os.Stdout)
	}
	fmt.Println(logMessage)
}

func getColorForLogLevel(level LogLevel) string {
	switch level {
	case InfoLevel:
		return "\033[37m" // Gray for info
	case DebugLevel:
		return "\033[34m" // Blue for debug
	case WarnLevel:
		return "\033[33m" // Yellow for warnings
	case ErrorLevel:
		return "\033[31m" // Red for errors
	case FatalLevel:
		return "\033[31;1m" // Bold red for fatal
	}
	return ""
}

// log first, then reset tags

func (l *Logger) Info(ctx context.Context, message any) *Logger {
	l.logInternal(ctx, InfoLevel, message)
	l.resetTags()
	return l
}

func (l *Logger) Warn(ctx context.Context, message any) *Logger {
	l.logInternal(ctx, WarnLevel, message)
	l.resetTags()
	return l
}

func (l *Logger) Error(ctx context.Context, message any) *Logger {
	l.logInternal(ctx, ErrorLevel, message)
	l.resetTags()
	return l
}

func (l *Logger) Debug(ctx context.Context, message any) *Logger {
	l.logInternal(ctx, DebugLevel, message)
	l.resetTags()
	return l
}

func (l *Logger) Fatal(ctx context.Context, message any) *Logger {
	l.logInternal(ctx, FatalLevel, message)
	l.resetTags()
	// Handle fatal error (e.g., exit the program)
	os.Exit(1)
	return l
}

// Add one or more slog.Attr to the logger
func (l *Logger) add(attrs ...slog.Attr) *Logger {
	// Lock the mutex before modifying the tags
	l.mu.Lock()
	defer l.mu.Unlock()

	// Add the attributes to the logger's tags for the current log
	l.tags = append(l.tags, attrs...)
	return l
}

// Reset the tags to avoid accumulation across log calls
func (l *Logger) resetTags() {
	// Lock the mutex before resetting the tags
	l.mu.Lock()
	defer l.mu.Unlock()

	l.tags = []slog.Attr{}
}
