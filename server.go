package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/kamsandhu93/gqldenring/logger"
	"github.com/kamsandhu93/gqldenring/middleware"
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
	log.SetFlags(log.Ldate | log.Ltime) //| log.Lshortfile)

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

	h := newHandler(resolver)

	http.Handle("/", middleware.WithReqID(logger.RequestIDKey, middleware.WithLogging(logger.LogID, playground.Handler("GraphQL playground", "/query"))))
	http.Handle("/query", h)
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

func newHandler(resolver *graph.Resolver) http.Handler {
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))
	srv.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		oc := graphql.GetOperationContext(ctx)
		logger.LogID(ctx, "[INFO] Incoming operation: %s %s %s", oc.OperationName, oc.Variables, strings.Replace(oc.RawQuery, "\n", " ", -1))
		return next(ctx)
	})
	srv.SetRecoverFunc(func(ctx context.Context, err interface{}) error {
		logger.LogID(ctx, "[ERROR] Panic caused by %v", err)
		debug.PrintStack()
		return gqlerror.Errorf("Internal server error!")
	})

	h := middleware.WithReqID(logger.RequestIDKey,
		middleware.WithLogging(logger.LogID,
			srv))

	return h
}
