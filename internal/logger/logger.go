package logger

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
)

// RequestIDKey Key to use when setting the request ID.
const RequestIDKey = "reqID"

// LogID is a thin wrapper around log that prefixes the request ID to the log msg.
func LogID(ctx context.Context, msg string, a ...any) {
	reqID := ctx.Value(RequestIDKey)

	v, ok := reqID.(string)
	if !ok {
		v = ""
	}
	_, path, lineNo, _ := runtime.Caller(1)
	file := filepath.Base(path)

	logPre := fmt.Sprintf("%s:%d %s %s", file, lineNo, v, msg)

	log.Printf(logPre, a...)
}
