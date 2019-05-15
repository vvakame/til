workflow "post draft of blog" {
  on = "push"
  resolves = ["Slack notification"]
}

action "filter PR merged" {
  uses = "actions/bin/filter@3c0b4f0e63ea54ea5df2914b4fabf383368cd0da"
  args = "merged true"
}

action "query PR contents" {
  uses = "Ilshidur/action-slack@f37693b4e0589604815219454efd5cb9b404fb85"
  needs = ["filter PR merged"]
}

action "Slack notification" {
  uses = "Ilshidur/action-slack@master"
  secrets = ["SLACK_WEBHOOK"]
  args = "A new commit has been pushed."
  needs = ["filter PR merged"]
}
