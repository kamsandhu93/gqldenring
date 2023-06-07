package logger

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

// RequestIDKey Key to use when setting the request ID.
const RequestIDKey = "reqID"

var ErrorLogger *log.Logger
var InfoLogger *log.Logger

func init() {
	InfoLogger = log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime)
	ErrorLogger = log.New(os.Stderr, "[ERROR]\t", log.Ldate|log.Ltime)
}

// Info is a thin wrapper around log that prefixes the request ID to the log msg.
func Info(ctx context.Context, msg string, a ...any) {
	reqID := ctx.Value(RequestIDKey)

	v, ok := reqID.(string)
	if !ok {
		v = ""

	}
	_, path, lineNo, _ := runtime.Caller(1)
	file := filepath.Base(path)

	logPre := fmt.Sprintf("%s:%d %s %s", file, lineNo, v, msg)

	InfoLogger.Printf(logPre, a...)
}

func Error(ctx context.Context, msg string, a ...any) {
	reqID := ctx.Value(RequestIDKey)

	v, ok := reqID.(string)
	if !ok {
		v = ""

	}
	_, path, lineNo, _ := runtime.Caller(1)
	file := filepath.Base(path)

	logPre := fmt.Sprintf("%s:%d %s %s", file, lineNo, v, msg)

	ErrorLogger.Printf(logPre, a...)
}
