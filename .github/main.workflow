workflow "make blog post" {
  resolves = [
    "cat GITHUB_EVENT_PATH",
    "blog to slack",
  ]
  on = "pull_request"
}

action "cat GITHUB_EVENT_PATH" {
  uses = "actions/bin/sh@master"
  args = ["cat $GITHUB_EVENT_PATH"]
}


action "filter PR merged" {
  uses = "actions/bin/filter@master"
  args = "merged true"
}

action "pr2md" {
  uses = "./github-actions/pr-to-md"
  secrets = ["GITHUB_TOKEN"]
  needs = ["cat"]
}

action "md2blog" {
  uses = "./github-actions/md-to-blogpost"
  args = ["--owner", "vvakame", "--name", "vvakame-blog", "--timezone", "Asia/Tokyo", "result.md"]
  secrets = ["BLOG_REPO_GITHUB_TOKEN"]
  needs = ["pr2md"]
}

action "blog to slack" {
  uses = "Ilshidur/action-slack@master"
  args = "PR merged!: {{ EVENT_PAYLOAD.pull_request.html_url }}"
  secrets = ["SLACK_WEBHOOK"]
  needs = ["md2blog"]
}
