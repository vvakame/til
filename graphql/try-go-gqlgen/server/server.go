package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/handler"
	trygogqlgen "github.com/vvakame/til/graphql/try-go-gqlgen"
	"github.com/vvakame/til/graphql/try-go-gqlgen/models"
	"fmt"
	"runtime/debug"
	"context"
	"github.com/pkg/errors"
	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/gqlerror"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	resolver := trygogqlgen.NewResolver()
	http.Handle("/", handler.Playground("GraphQL playground", "/query"))
	http.Handle("/query",
		models.DataloaderMiddleware(
			resolver.UserMap,
			handler.GraphQL(
				trygogqlgen.NewExecutableSchema(trygogqlgen.Config{
					Resolvers: resolver,
				}),
				handler.ErrorPresenter(func(ctx context.Context, err error) *gqlerror.Error {
					err = errors.Cause(err)

					rc := graphql.GetResolverContext(ctx)

					return &gqlerror.Error{
						Message: err.Error(),
						Path:    rc.Path,
						Extensions: map[string]interface{}{
							"hello": "myError!",
						},
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

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
