package middleware

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
)

var prefix string
var reqid uint64

func init() {
	hostname, err := os.Hostname()
	if hostname == "" || err != nil {
		hostname = "localhost"
	}
	var buf [12]byte
	var b64 string
	for len(b64) < 10 {
		_, err := rand.Read(buf[:])
		if err != nil {
			panic("Cant generate random int for process identifier")
		}
		b64 = base64.StdEncoding.EncodeToString(buf[:])
		b64 = strings.NewReplacer("+", "", "/", "").Replace(b64)
	}

	prefix = fmt.Sprintf("%s/%s", hostname, b64[0:10])
}

// WithReqID is a middleware that injects a request ID into the context of each
// request. A request ID is a string of the form "host.example.com/random-0001",
// where "random" is a base62 random string that uniquely identifies this go
// process, and where the last number is an atomically incremented request
// counter.
func WithReqID(ctxReqIDKey string, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		newReqID := atomic.AddUint64(&reqid, 1)
		//nolint:staticcheck
		ctx = context.WithValue(ctx, ctxReqIDKey, fmt.Sprintf("%s-%06d", prefix, newReqID))

		r = r.WithContext(ctx)

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
