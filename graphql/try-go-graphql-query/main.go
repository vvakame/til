package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/k0kubun/pp"
	"github.com/motemen/go-graphql-query"
)

func main() {
	type Query struct {
		User *struct {
			GraphQLArguments struct {
				Login string `graphql:"$login,notnull"`
			}
			Name          string
			Organizations *struct {
				GraphQLArguments struct {
					First int `graphql:"$organizationFirst,notnull"`
				}
				Edges []*struct {
					Cursor string // Edge みたいなembed structの展開は思い通りじゃなかった
					Node   *struct {
						Name string
					}
				}
				Nodes []*struct {
					Name string
				}
			} `graphql:"(first: $organizationFirst)"`
		} `graphql:"(login: $login)"`

		TilRepository *struct {
			GraphQLArguments struct {
				Owner string `graphql:"$tilOwner,notnull"`
				Name  string `graphql:"$tilRepo,notnull"`
			}
			Name string
		} `graphql:"alias=repository,(owner: $tilOwner, name: $tilRepo)"`

		TsfmtRepository *struct {
			GraphQLArguments struct {
				Owner string `graphql:"$tsfmtOwner,notnull"`
				Name  string `graphql:"$tsfmtRepo,notnull"`
			}
			Name string
		} `graphql:"alias=repository,(owner: $tsfmtOwner, name: $tsfmtRepo)"`
	}

	obj := &Query{}

	query, err := graphqlquery.Build(obj)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(query))

	req := &struct {
		Query     string                 `json:"query"`
		Variables map[string]interface{} `json:"variables"`
	}{
		Query: string(query),
		Variables: map[string]interface{}{
			"login":             "vvakame",
			"organizationFirst": 100,
			"tilOwner":          "vvakame",
			"tilRepo":           "til",
			"tsfmtOwner":        "vvakame",
			"tsfmtRepo":         "typescript-formatter",
		},
	}
	reqB, err := json.MarshalIndent(req, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(reqB))

	apiReq, err := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewReader(reqB))
	if err != nil {
		panic(err)
	}
	apiReq.Header.Add("Content-Type", "application/json")
	apiReq.Header.Add("Authorization", "bearer "+os.Getenv("GITHUB_TOKEN"))
	resp, err := http.DefaultClient.Do(apiReq)
	if err != nil {
		panic(err)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))

	type Resp struct {
		Data *Query `json:"data"`
	}
	err = json.Unmarshal(b, &Resp{obj})
	if err != nil {
		panic(err)
	}

	pp.Println(obj)
}
