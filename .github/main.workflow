workflow "post draft of blog" {
  resolves = [
    "Slack notification",
  ]
  on = "pull_request"
}

action "filter PR merged" {
  uses = "actions/bin/filter@3c0b4f0e63ea54ea5df2914b4fabf383368cd0da"
  args = "merged true"
  needs = ["cat pr2md"]
}

action "Slack notification" {
  uses = "Ilshidur/action-slack@master"
  secrets = ["SLACK_WEBHOOK"]
  args = "A new commit has been pushed."
  needs = ["filter PR merged"]
}

action "cat" {
  uses = "actions/bin/sh@master"
  args = ["cat $GITHUB_EVENT_PATH"]
}

action "ls" {
  uses = "actions/bin/sh@master"
  args = ["ls -ltr"]
  needs = ["cat"]
}

action "pr2md" {
  uses = "./github-actions/pr-to-md"
  needs = ["ls"]
  secrets = ["BLOG_REPO_GITHUB_TOKEN"]
}

action "cat pr2md" {
  uses = "actions/bin/sh@master"
  args = ["cat result.md"]
  needs = ["pr2md"]
}
