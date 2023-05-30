package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/kamsandhu93/gqldenring/graph"
	"github.com/kamsandhu93/gqldenring/graph/generated"
)

const defaultPort = "8080"

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	srv.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		oc := graphql.GetOperationContext(ctx)
		log.Printf("[INFO] Incoming operation: %s %s", oc.OperationName, strings.Replace(oc.RawQuery, "\n", " ", -1))
		return next(ctx)
	})

	http.Handle("/", withLogging(playground.Handler("GraphQL playground", "/query")))
	http.Handle("/query", withLogging(srv))
	http.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
		_, err := io.WriteString(writer, "okay\n")
		if err != nil {
			panic(1)
		}
	})

	log.Printf("Starting server version=%s commit=%s date=%s connect to http://localhost:%s/ for GraphQL playground",
		version, commit, date, port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func withLogging(h http.Handler) http.Handler {
	logFn := func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()

		uri := r.RequestURI
		method := r.Method
		h.ServeHTTP(rw, r) // serve the original request

		duration := time.Since(start)
		log.Printf("%s %s %s", uri, method, duration)
	}
	return http.HandlerFunc(logFn)
}
