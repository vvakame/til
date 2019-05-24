package main

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"path"
	"time"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
)

var (
	githubToken  = kingpin.Flag("github_token", "GitHub token for GitHub endpoint request.").Default(os.Getenv("GITHUB_TOKEN")).String()
	owner        = kingpin.Flag("owner", "name of repository owner.").Required().String()
	name         = kingpin.Flag("name", "name of repository.").Required().String()
	prNumber     = kingpin.Flag("pr_number", "number of pull request.").Required().Int()
	templatePath = kingpin.Flag("template_path", "markdown template path. it uses with html/template").String()
	outputPath   = kingpin.Arg("output_path", "result output path").Default("result.md").String()
)

func main() {
	kingpin.Parse()

	ctx := context.Background()

	if *githubToken == "" {
		log.Fatal("--github_token or $GITHUB_TOKEN is required")
	}

	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: *githubToken})
	httpClient := oauth2.NewClient(ctx, src)

	client := githubv4.NewClient(httpClient)
	resp, err := getPRInfo(ctx, client, *owner, *name, *prNumber)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.OpenFile(*outputPath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}

	err = generateMarkdown(ctx, f, *templatePath, resp)
	if err != nil {
		log.Fatal(err)
	}
}

func generateMarkdown(ctx context.Context, w io.Writer, templatePath string, resp *PRInfo) error {

	var templateName string
	if templatePath == "" {
		templatePath = "./markdown.tmpl.md"
		templateName = "markdown.tmpl.md"
	} else {
		templateName = path.Base(templatePath)
	}

	pr := resp.Repository.PullRequest

	_, _ = fmt.Fprintf(w, "---\n")
	err := yaml.NewEncoder(w).Encode(struct {
		Title string   `yaml:"title"`
		Date  string   `yaml:"date"`
		Tags  []string `yaml:"tags"`
	}{
		Title: string(pr.Title),
		Date:  time.Now().Format("2006-01-02 15:04:05"), // TODO Timezone
		Tags:  []string{},                               // TODO
	})
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintf(w, "---\n\n")

	samePrevStock := make(map[string]string)
	tmpl, err := template.
		New(templateName).
		Funcs(map[string]interface{}{
			"date": func(t githubv4.DateTime) string {
				return t.Format("2006-01-02 15:04:05") // TODO Timezone
			},
			"isSamePrev": func(group string, value githubv4.String) bool {
				if samePrevStock[group] == string(value) {
					return true
				}
				samePrevStock[group] = string(value)
				return false
			},
		}).
		ParseFiles(templatePath)
	if err != nil {
		return err
	}

	return tmpl.Execute(w, &pr)
}

type PRInfo struct {
	Repository struct {
		PullRequest struct {
			Title  githubv4.String
			URL    githubv4.String
			Body   githubv4.String
			Author struct {
				Login     githubv4.String
				AvatarURL githubv4.String `graphql:"avatarUrl(size: $avatarImageSize)"`
			}
			CreatedAt githubv4.DateTime
			Files     struct {
				PageInfo struct {
					HasNextPage githubv4.Boolean
					EndCursor   githubv4.String
				}
				TotalCount githubv4.Int
				Nodes      []struct {
					Path      githubv4.String
					Additions githubv4.Int
					Deletions githubv4.Int
				}
			} `graphql:"files(first: 100)"`
			Comments struct {
				PageInfo struct {
					HasNextPage githubv4.Boolean
					EndCursor   githubv4.String
				}
				TotalCount githubv4.Int
				Nodes      []struct {
					Author struct {
						Login     githubv4.String
						AvatarURL githubv4.String `graphql:"avatarUrl(size: $avatarImageSize)"`
					}
					Body      githubv4.String
					CreatedAt githubv4.DateTime
				}
			} `graphql:"comments(first: 100)"`
		} `graphql:"pullRequest(number: $prNumber)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

func getPRInfo(ctx context.Context, cli *githubv4.Client, owner, name string, prNumber int) (*PRInfo, error) {
	var query PRInfo
	variables := map[string]interface{}{
		"owner":           githubv4.String(owner),
		"name":            githubv4.String(name),
		"prNumber":        githubv4.Int(prNumber),
		"avatarImageSize": githubv4.Int(64),
	}

	err := cli.Query(ctx, &query, variables)
	return &query, err
}
