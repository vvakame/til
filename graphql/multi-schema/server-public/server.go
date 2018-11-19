package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/handler"
	api "github.com/vvakame/til/graphql/multi-schema/api-public"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	http.Handle("/", handler.Playground("GraphQL playground", "/query"))
	http.Handle("/query",
		handler.GraphQL(
			api.NewExecutableSchema(
				api.Config{
					Resolvers: api.NewResolver(),
				},
			),
		),
	)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
