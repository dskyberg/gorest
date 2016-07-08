base =  <<EOF
github accepts the following commands:
- new: create a new issue
- get: display an issue
- update: Update an issue
- close: Mark an issue as closed.
- help: this text
EOF

new = <<EOF
github *new* creates a new issue in a GitHub repository

Here are the supported keys:

* title:      Required. Issue title
* body:       Optional. Text to add to the issue body
* labels:     Optional A comma separated list of labels
* milestone:  Optional. The milestone must exist in the identified repository
* assignee:   Optional.  Github user to assign the issue to
* repo:       Optional. If not provided, the default repo will be used.

Do not use the '=' character within a value.
Unrecognized keys are ignored.

Examples:
* /github new title = Here is my title  labels = label1, label2
* /github new title = Here is my title  labels = label1, label2 repo = my repo
EOF

get = <<EOF
github *get* fetches an issue and displays it.

Here are the supported keys:
- number:      Required. Issue Number.  Either provided after "get", or as a KV
- repo:        Optional. If not provided, the default repo will be used.

Do not use the '=' character within a value.
Unrecognized keys are ignored.

Examples:
  /github get 152
  /github get 152 repo = my repo
  /github get repo= my repo number = 152
EOF

// This is here for testing.  Do not remove.
test.some.stuff = "Help for test some stuff"
