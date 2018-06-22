package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"

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
				handler.RecoverFunc(func(ctx context.Context, err interface{}) error {
					// panicからの回復方法
					fmt.Fprintln(os.Stderr, err)
					fmt.Fprintln(os.Stderr)
					debug.PrintStack()

					return nil
				}),
				handler.ResolverMiddleware(func(ctx context.Context, next graphql.Resolver) (interface{}, error) {
					// Resolverが1回仕事することになるたびに呼ばれるようだ
					res, err := next(ctx)
					return res, err
				}),
				handler.RequestMiddleware(func(ctx context.Context, next func(ctx context.Context) []byte) []byte {
					// 結果のJSONの "data" のvalue部分のbyte列が帰ってくるようだ
					b := next(ctx)
					return b
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
