//go:generate statik -src=./misc

package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/google/go-github/v25/github"
	"golang.org/x/oauth2"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	githubToken = kingpin.Flag("github_token", "GitHub token for GitHub endpoint request.").Default(os.Getenv("BLOG_REPO_GITHUB_TOKEN")).String()

	githubEventPath = kingpin.Flag("github_event_path", "GitHub event data json path.").Default(os.Getenv("GITHUB_EVENT_PATH")).String()

	owner = kingpin.Flag("owner", "name of repository owner.").String()
	name  = kingpin.Flag("name", "name of repository.").Required().String()

	baseBranch = kingpin.Flag("base_branch", "name of base branch.").Default("master").String()

	commitMessage = kingpin.Flag("commit_message", "commit message text.").Default("auto generated commit").String()
	timezone      = kingpin.Flag("timezone", "timezone for blog post date").Default("UTC").String()
	postPath      = kingpin.Flag("post_path", "blog post markdown path in repository.").String()
	imagePath     = kingpin.Flag("image_path", "blog post image path in repository.").String()
	baseImageURL  = kingpin.Flag("base_image_url", "base image url in blog site.").Default("/images").String()

	prBranch = kingpin.Flag("pr_branch", "name of pull request branch.").String()
	prTitle  = kingpin.Flag("pr_title", "title of pull request.").String()
	prBody   = kingpin.Flag("pr_body", "body of pull request.").String()

	inputPath = kingpin.Arg("input_path", "input markdown path").Required().String()
)

func main() {
	kingpin.Parse()

	ctx := context.Background()

	if *githubToken == "" {
		log.Fatal("--github_token or $GITHUB_TOKEN is required")
	}

	eventData := &GitHubEvent{}
	{
		b, err := ioutil.ReadFile(*githubEventPath)
		if os.IsNotExist(err) {
			// ok
		} else if err != nil {
			log.Fatal(err)
		} else {
			err = json.Unmarshal(b, eventData)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	if *owner == "" {
		*owner = eventData.Repository.Owner.Login
	}
	if *postPath == "" {
		loc, err := time.LoadLocation(*timezone)
		if err != nil {
			log.Fatal(err)
		}
		var postDate string
		if eventData.PullRequest.MergedAt.IsZero() {
			postDate = eventData.PullRequest.MergedAt.In(loc).Format("2006-01-02")
		} else {
			postDate = time.Now().In(loc).Format("2006-01-02")
		}
		*postPath = fmt.Sprintf("source/_posts/%s-%s.md", postDate, eventData.PullRequest.Head.Ref)
	}
	if *imagePath == "" {
		loc, err := time.LoadLocation(*timezone)
		if err != nil {
			log.Fatal(err)
		}
		var postDate string
		if eventData.PullRequest.MergedAt.IsZero() {
			postDate = eventData.PullRequest.MergedAt.In(loc).Format("2006-01-02")
		} else {
			postDate = time.Now().In(loc).Format("2006-01-02")
		}
		*imagePath = fmt.Sprintf("source/images/%s-%s", postDate, eventData.PullRequest.Head.Ref)
	}
	if *prBranch == "" {
		*prBranch = fmt.Sprintf("from-pr-%d", eventData.Number)
	}
	if *prTitle == "" {
		*prTitle = fmt.Sprintf("blog post from '%s'", eventData.PullRequest.Title)
	}
	if *prBody == "" {
		*prBody = fmt.Sprintf("from %s", eventData.PullRequest.HTMLURL)
	}

	client := github.NewClient(
		oauth2.NewClient(
			ctx,
			oauth2.StaticTokenSource(&oauth2.Token{AccessToken: *githubToken}),
		),
	)

	contentMap, err := createContent(ctx, &CreateContentReq{
		InputPath:    *inputPath,
		PostPath:     *postPath,
		ImagePath:    *imagePath,
		BaseImageURL: *baseImageURL,
	})
	if err != nil {
		log.Fatal(err)
	}

	branchRef, commit, err := createNewBranch(ctx, &CreateCommitReq{
		Client: client,
		Owner:  *owner,
		Name:   *name,

		BaseRef:   fmt.Sprintf("refs/heads/%s", *baseBranch),
		BranchRef: fmt.Sprintf("refs/heads/%s", *prBranch),

		Message: *commitMessage,
		Files:   contentMap,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println(branchRef.GetURL())
	log.Println(commit.GetURL())

	pr, err := createPR(ctx, &CreatePRReq{
		Client: client,
		Owner:  *owner,
		Name:   *name,
		Title:  *prTitle,
		Body:   *prBody,
		Base:   *baseBranch,
		Head:   *prBranch,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println(pr.GetURL())
}

type CreateContentReq struct {
	InputPath string

	PostPath     string
	ImagePath    string
	BaseImageURL string
}

func createContent(ctx context.Context, req *CreateContentReq) (map[string][]byte, error) {

	inputPath := req.InputPath
	postPath := req.PostPath
	imagePath := req.ImagePath
	baseImageURL := req.BaseImageURL

	b, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return nil, err
	}

	// NOTE
	//   AST → markdown できるいい感じのライブラリないかな…
	//   🙅 https://github.com/russross/blackfriday
	//   🙅 https://godoc.org/github.com/gomarkdown/markdown

	re := regexp.MustCompile(`!\[(?P<alt>[^]]*)]\((?P<url>[^)]+)\)`)

	getContent := func(imageURL string) ([]byte, error) {
		parsed, err := url.Parse(imageURL)
		if err != nil {
			return nil, err
		}

		if parsed.Host == "" {
			return ioutil.ReadFile(imageURL)
		}

		resp, err := http.Get(imageURL)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode < 200 && 300 <= resp.StatusCode {
			return nil, fmt.Errorf("fetch image failed: %d %s", resp.StatusCode, imageURL)
		}
		return ioutil.ReadAll(resp.Body)
	}

	contentMap := make(map[string][]byte)

	var buf bytes.Buffer
	tail := 0
	for _, submatches := range re.FindAllSubmatchIndex(b, -1) {
		alt := string(b[submatches[2]:submatches[3]])
		imageURLStr := string(b[submatches[4]:submatches[5]])
		imageURL, err := url.Parse(imageURLStr)
		if err != nil {
			return nil, err
		}
		shouldCopy := true
		if imageURL.Host == "github.com" {
			// GitHubのユーザアイコン部分はコピーせず直リンクにする
			// https://github.com/vvakame.png?size=64 など
			shouldCopy = false
		}

		imageBlob, err := getContent(imageURLStr)
		if err != nil {
			return nil, err
		}
		mimeType := http.DetectContentType(imageBlob)
		ext := strings.TrimPrefix(mimeType, "image/")
		fileName := fmt.Sprintf("%x.%s", md5.Sum(imageBlob), ext)

		// before matching
		buf.Write(b[tail:submatches[0]])
		// ![
		buf.Write(b[submatches[0]:submatches[2]])
		// alt
		buf.WriteString(alt)
		// ](
		buf.Write(b[submatches[3]:submatches[4]])
		// imageURL
		if shouldCopy {
			buf.WriteString(path.Join(baseImageURL, fileName))
		} else {
			buf.WriteString(imageURLStr)
		}
		// )
		buf.Write(b[submatches[5]:submatches[1]])

		tail = submatches[1]

		if shouldCopy {
			contentMap[path.Join(imagePath, fileName)] = imageBlob
		}
	}
	buf.Write(b[tail:])

	contentMap[postPath] = buf.Bytes()

	return contentMap, nil
}

type CreateCommitReq struct {
	Client *github.Client
	Owner  string
	Name   string

	BranchRef string

	Message string
	BaseRef string
	Files   map[string][]byte
}

func createNewBranch(ctx context.Context, req *CreateCommitReq) (*github.Reference, *github.Commit, error) {

	client := req.Client
	owner := req.Owner
	name := req.Name

	baseRef := req.BaseRef
	if baseRef == "" {
		baseRef = "refs/heads/master"
	}

	stringPtr := func(s string) *string {
		return &s
	}

	ref, _, err := client.Git.GetRef(ctx, owner, name, baseRef)
	if err != nil {
		return nil, nil, err
	}

	commit, _, err := client.Git.GetCommit(ctx, owner, name, ref.GetObject().GetSHA())
	if err != nil {
		return nil, nil, err
	}

	var entries []github.TreeEntry
	for filePath, content := range req.Files {
		blob, _, err := client.Git.CreateBlob(ctx, owner, name, &github.Blob{
			Encoding: stringPtr("base64"),
			Content:  stringPtr(base64.StdEncoding.EncodeToString(content)),
		})
		if err != nil {
			return nil, nil, err
		}

		entries = append(entries, github.TreeEntry{
			Path: stringPtr(filePath),
			Mode: stringPtr("100644"),
			Type: stringPtr("blob"),
			SHA:  stringPtr(blob.GetSHA()),
		})
	}

	tree, _, err := client.Git.CreateTree(ctx, owner, name, commit.GetTree().GetSHA(), entries)
	if err != nil {
		return nil, nil, err
	}

	newCommit, _, err := client.Git.CreateCommit(ctx, owner, name, &github.Commit{
		Message: &req.Message,
		Parents: []github.Commit{*commit},
		Tree:    tree,
	})
	if err != nil {
		return nil, nil, err
	}

	branchRef, _, err := client.Git.CreateRef(ctx, owner, name, &github.Reference{
		Ref: stringPtr(req.BranchRef),
		Object: &github.GitObject{
			SHA: stringPtr(newCommit.GetSHA()),
		},
	})
	if err != nil {
		return nil, nil, err
	}

	return branchRef, newCommit, nil
}

type CreatePRReq struct {
	Client *github.Client
	Owner  string
	Name   string

	Title string
	Body  string
	Base  string
	Head  string
}

func createPR(ctx context.Context, req *CreatePRReq) (*github.PullRequest, error) {

	client := req.Client
	owner := req.Owner
	name := req.Name

	stringPtr := func(s string) *string {
		return &s
	}

	if req.Base == "" {
		req.Base = "master"
	}

	maintainerCanModify := true

	pr, _, err := client.PullRequests.Create(ctx, owner, name, &github.NewPullRequest{
		Title:               stringPtr(req.Title),
		Head:                stringPtr(req.Head),
		Base:                stringPtr(req.Base),
		Body:                stringPtr(req.Body),
		MaintainerCanModify: &maintainerCanModify,
	})
	if err != nil {
		return nil, err
	}

	return pr, nil
}

type GitHubEvent struct {
	Number      int
	PullRequest struct {
		Title    string
		HTMLURL  string    `json:"html_url"`
		MergedAt time.Time `json:"merged_at"`
		Head     struct {
			Ref string
		}
	} `json:"pull_request"`
	Repository struct {
		Name  string
		Owner struct {
			Login string
		}
	}
}
