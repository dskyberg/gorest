package commands

import (
  . "github.com/confyrm/gorest/slack/command"
)

// This is the list of commands that is given to SlashRouter.
// Each command is added according to the slash command that will be
// presented by Slack.  The command handler must be in this "commands" package.
var SlashCommands = Commands {
  Command{
    "/devhub",
    false,
    DevHub,
  },
}
