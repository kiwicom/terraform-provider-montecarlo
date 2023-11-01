resource "montecarlo_iam_member" "example_builtin" {
  group = "groups/editors-all"
  member = "user:user@google.com"
}

resource "montecarlo_iam_member" "example_custom" {
  group = "groups/custom-group"
  member = "user:user@google.com"
}

resource "montecarlo_iam_member" "example_multiple" {
  group = "groups/custom-group"
  member = each.value
  for_each = toset([
    "user:user1@google.com",
    "user:user2@google.com"
  ])
}
