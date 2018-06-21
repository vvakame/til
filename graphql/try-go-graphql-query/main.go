package main

import (
	"fmt"

	"github.com/motemen/go-graphql-query"
)

func main() {
	type Edge struct {
		Cursor string
	}

	type Organization struct {
		Name string
	}

	type OrganizationEdge struct {
		Cursor string // Edge みたいなembed structの展開は思い通りじゃなかった
		Node   *Organization
	}

	type OrganizationConnection struct {
		Edges []*OrganizationEdge
		Nodes []*Organization
	}

	type User struct {
		Name          string
		Organizations []*OrganizationConnection `graphql:"(first: $organizationFirst)"`
	}

	type Repository struct {
		Name string
	}

	type Test struct {
		GraphQLArguments struct {
			User              string `graphql:"$login,notnull"`
			OrganizationFirst int    `graphql:"$organizationFirst,notnull"`
		}
		User            *User       `graphql:"(login: $login)"`
		TilRepository   *Repository `graphql:"alias=repository,(owner: \"vvakame\", name: \"til\")"`
		TsfmtRepository *Repository `graphql:"alias=repository,(owner: \"vvakame\", name: \"typescript-formatter\")"`
	}

	obj := &Test{}

	query, err := graphqlquery.Build(obj)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(query))
}
