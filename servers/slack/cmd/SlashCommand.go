package cmd

import (
  "github.com/confyrm/gorest/slack"
  . "github.com/confyrm/gorest/errors"
  "github.com/confyrm/gorest/config"
)

// CommandHandlerFunc is mapped to a slash command. The function returns the
// outcome as a string, or it returns a StatusError. If the CommandHandlerFunc
// is not long running, then the StatusError is returned by the route handler.
// If the CommandHandlerFunc is long running, then it is simply logged.
type SlashCommandHandlerFunc func(config *config.Config, sReq *slack.Request) (*slack.Response, *StatusError)

// Command, like Route, maps slash commands with command handlers
type SlashCommand struct {

  // The slash command.  Such as /devhub
  Name string

  // If this is a long command, it will be run as a goroutine.
  // The request MUST contain a response_url, or an error will be thrown
  IsLong bool

  // The handler.
  Handler SlashCommandHandlerFunc
}
// Helper type.
type SlashCommands []SlashCommand

// Router to map commands and handler funcs.
type SlashCommandRouter map[string]SlashCommand

// New turns the set of Commands into a map for lookup
func (cmds SlashCommands) New() map[string]SlashCommand {
  m := make(map[string]SlashCommand)
  for _, cmd := range cmds {
    m[cmd.Name] = cmd
  }
  return m
}
