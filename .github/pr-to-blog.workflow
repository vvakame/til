workflow "post draft of blog" {
  resolves = [
    "Slack notification",
  ]
  on = "pull_request"
}

action "filter PR merged" {
  uses = "actions/bin/filter@3c0b4f0e63ea54ea5df2914b4fabf383368cd0da"
  args = "merged true"
  needs = ["pr2md"]
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

action "pr2md" {
  uses = "vvakame/til/github-actions/pr-to-md@kick-actions-r4"
  args = ["cat $GITHUB_EVENT_PATH"]
  needs = ["cat"]
}
