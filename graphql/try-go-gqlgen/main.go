package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/vektah/gqlgen/handler"
	"github.com/vvakame/til/graphql/try-go-gqlgen/graph"
	"github.com/vvakame/til/graphql/try-go-gqlgen/models"
)

func main() {
	app := graph.NewMyApp()
	http.Handle("/", handler.Playground("Todo", "/query"))
	http.Handle("/query",
		models.DataloaderMiddleware(
			app.UserMap,
			handler.GraphQL(graph.MakeExecutableSchema(app)),
		),
	)

	fmt.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
