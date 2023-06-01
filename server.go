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
	"github.com/vektah/gqlparser/v2/gqlerror"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	memDB "github.com/kamsandhu93/gqldenring/db/memory"
	sqlDB "github.com/kamsandhu93/gqldenring/db/sql"
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
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	sqlConn := os.Getenv("SQL_CONN")

	var db graph.DB
	if sqlConn != "" {
		log.Print("[INFO] Using sql db")
		db = sqlDB.NewDB(sqlConn)
	} else {
		log.Print("[INFO] Using in memory db")
		db = memDB.NewDB()
	}

	resolver := graph.NewResolver(db)

	srv := newServer(resolver)

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

func newServer(resolver *graph.Resolver) *handler.Server {
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))
	srv.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		oc := graphql.GetOperationContext(ctx)
		log.Printf("[INFO] Incoming operation: %s %s %s", oc.OperationName, oc.Variables, strings.Replace(oc.RawQuery, "\n", " ", -1))
		return next(ctx)
	})
	srv.SetRecoverFunc(func(ctx context.Context, err interface{}) error {
		log.Printf("[ERROR] Panic caused by %v", err)

		return gqlerror.Errorf("Internal server error!")
	})
	return srv
}

func withLogging(h http.Handler) http.Handler {
	logFn := func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()

		uri := r.RequestURI
		method := r.Method
		h.ServeHTTP(rw, r) // serve the original request

		duration := time.Since(start)
		log.Printf("[INFO] Request complete %s %s %s", uri, method, duration)
	}
	return http.HandlerFunc(logFn)
}
