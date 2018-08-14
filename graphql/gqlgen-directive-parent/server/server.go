package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/handler"
	"github.com/vvakame/til/graphql/gqlgen-directive-parent"
)

const defaultPort = "8080"

func CurrentUser(ctx context.Context) *gqlgen_directive_parent.User {
	return &gqlgen_directive_parent.User{
		ID:   "User:123",
		Name: "123-san",
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	http.Handle("/", handler.Playground("GraphQL playground", "/query"))
	http.Handle("/query",
		handler.GraphQL(
			gqlgen_directive_parent.NewExecutableSchema(
				gqlgen_directive_parent.Config{
					Resolvers: gqlgen_directive_parent.NewResolver(),
					Directives: gqlgen_directive_parent.DirectiveRoot{
						LoginUser: func(ctx context.Context, next graphql.Resolver) (interface{}, error) {

							user := CurrentUser(ctx)
							if user == nil {
								return nil, nil
							}

							resp, err := next(ctx)
							if err != nil {
								return nil, err
							}

							// resp is string from Todo#text.
							fmt.Println(resp)

							// How can I get the parent object?
							rctx := graphql.GetResolverContext(ctx)
							fmt.Println(rctx)
							// I need Todo object. but How?
							var todo gqlgen_directive_parent.Todo
							if todo.User.ID != user.ID {
								return nil, nil
							}

							return resp, nil
						},
					},
				},
			),
		),
	)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
