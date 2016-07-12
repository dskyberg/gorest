package routes

import (
	"path/filepath"
  "errors"
  "net/http"
  "fmt"
  "log"
	"strings"
  "encoding/json"

  "github.com/confyrm/gorest/slack"
  . "github.com/confyrm/gorest/errors"
  "github.com/confyrm/gorest/config"
    "github.com/confyrm/gorest/help"
  . "github.com/confyrm/gorest/servers/slack/commands"
)



// Make a static decoder for performance

var commandRouter = SlashCommands.New()

// This is set the first time HelpResponse is called.  It is set by Loading
// The help text file located in ../../../help/help.hcl
var helpResponses help.Help

// SlashRouter is the top level slash command router.
func SlashRouter(config *config.Config, rw http.ResponseWriter, req *http.Request) error {

  sReq := &slack.Request{}

  if err := sReq.DecodeHttp(req); err != nil {
    return StatusError{http.StatusBadRequest,  err}
  }
  // Dump the slack.Request to the log
  sReq.Log()

  command, err := sReq.TextToCommand()
  if err != nil {
    return StatusError{http.StatusInternalServerError, err}
  }
  log.Printf("Received Slack slash command: %#v", command)

  // First authZ test.  Ensure the token provided in the slack request
  // matches the SLACK_TOKEN
  if sReq.Token != config.GetString("SLACK_TOKEN") {
    return StatusError{http.StatusUnauthorized,
      errors.New("Not authorized. Wrong Slack Token.")}
  }

	// Look to see if we just need to return some help
	yes, err := HadHelp(config, command, rw)
	if err != nil {
		return StatusError{http.StatusInternalServerError, err}
	}
	if yes {
		// HadHelp send a Help response.  We're done.
		return nil
	}

  // Get the route handler for this slash command, such as '/devhub'
  route := commandRouter.Route(sReq.Command)
  if route == nil {
    // Oops!  No route found.  Must be an unknown command
    return StatusError {
      http.StatusBadRequest,
      fmt.Errorf("Command not found [%s]", sReq.Command),
    }
  }

  // At this point, we can either process the command and return a
  // slack.Response, or, we can kick off a goroutine, and return a
  // quick, happy response.  The long running command can send its
  //  response as a POST request to the response_url listed in the sReq.
  var response *slack.Response
  if route.IsLong {
    // Long running command. Kick off a goroutine and return a happy response.
    go route.Handler(config, sReq, command)
    // Create a quick response to let the user know the comand is running
    response = HappyResponse(command)
  } else {
    // Short short command.  Just run and return the response.
    var statusErr error
    response, statusErr = route.Handler(config, sReq, command)
    if statusErr != nil {
      return statusErr
    }
  }

  // The response is either the returned response from a short running command
  // handler, or a "command in process" style response.
  rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
  if err := json.NewEncoder(rw).Encode(response); err != nil {
    return StatusError{http.StatusInternalServerError,
      errors.New(fmt.Sprintf("Failure while encoding response data: %#v", err.Error()))}
  }
  return nil

}

func HadHelp(config *config.Config, command *slack.DevHubCommand, rw http.ResponseWriter) (bool, error) {

	if helpPath, ok := command.HelpPath(); ok {
    rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
    response := HelpResponse(config, helpPath)
    if err := json.NewEncoder(rw).Encode(response); err != nil {
      return true, StatusError{http.StatusInternalServerError, err}
    }
    return true, nil
  }

	// There was no help to process.
	return false, nil
}

// HappyResponse just sends back a "message received" response.  This is
// Sent when the command is long running, to let the user know that it's running
func HappyResponse(command *slack.DevHubCommand) *slack.Response {
  var cmdText string
  if len(command.Commands) > 0 {
    cmdText = command.Commands[0]
  } else {
    cmdText = ""
  }
  text := fmt.Sprintf("Roger that!  Message received!\r\nYour %s request is in process!", cmdText)
  response := slack.Response{slack.Ephemeral.String(), &text, nil}
  return &response
}

// HelpResponse looks up help text from the global helpResponses variable.
// The slice of commands are joined to create the lookup key.
func HelpResponse(config *config.Config, commands slack.Commands) *slack.Response {

  // Although called every time, this will only attempted to read the help
  // text the first time its called.
  ImportHelpText(config)

	cLen := len(commands)
  var text string
  // If the only command is help, then return the top level help
  if (cLen == 0) {
    text = helpResponses.Base()
  } else {
    // Create the help lookup key, by joining all the commands
    key := strings.Join(commands, help.Sep)
    text = helpResponses.Get(key)
  }

  response := slack.Response{slack.Ephemeral.String(), &text, nil}
  return &response
}

// ImportHelpText attempts to read an HCL formated file located in
// APP_ROOT/help/help.hcl.
func ImportHelpText(config *config.Config) {
  if helpResponses != nil {
    return
  }

  helpPath := filepath.Join(config.GetString("APP_ROOT"), "help.hcl" )
  log.Printf("Reading help from %s", helpPath)
  var err error
  helpResponses, err = help.ParseHelpFile(helpPath)
  if err != nil {
    log.Printf("Error parsing help file: ", err.Error())
  }
}
