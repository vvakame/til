package main

import (
	"fmt"
	"log"
	"net/http"

	"context"

	"github.com/pkg/errors"
	"github.com/vektah/gqlgen/graphql"
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
			handler.GraphQL(
				graph.MakeExecutableSchema(app),
				handler.ErrorPresenter(func(ctx context.Context, err error) error {
					err = errors.Cause(err)

					rc := graphql.GetResolverContext(ctx)

					return &MyError{
						graphql.ResolverError{
							Message: err.Error(),
							Path:    rc.Path,
						},
						"foobar!",
					}
				}),
			),
		),
	)

	fmt.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type MyError struct {
	graphql.ResolverError
	FooBar string `json:"foobar"`
}

func (e *MyError) Error() string {
	return e.ResolverError.Error()
}
