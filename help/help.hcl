base =  <<EOF
github accepts the following commands:
- new: create a new issue
- get: display an issue
- update: Update an issue
- close: Mark an issue as closed.
- help: this text
EOF

new = <<EOF
github new accepts a set of <key>=<value> statements.  Do not use
the '=' character within a value.  Unrecognized keys are ignored.
Here are the supported keys:
---
- title:      Required. Issue title
- body:       Optional. Text to add to the issue body
- labels:     Optionall A comma separated list of labels
- milestone:  Optional. The milestone must exist in the identified repository
- assignee:   Optional.  Github user to assign the issue to
- repository: Optional. If not provided, devhub will be used.
---
Example: /github new title = Here is my title labels = label1, label2
EOF

// This is here for testing.  Do not remove.
test.some.stuff = "Help for test some stuff"
