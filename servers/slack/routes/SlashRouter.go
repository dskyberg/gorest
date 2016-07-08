package routes

import (
  "io/ioutil"
	"path/filepath"
  "errors"
  "net/http"
  "fmt"
  "log"
  "strings"
  "encoding/json"
  "github.com/gorilla/schema"
  "github.com/hashicorp/hcl"

  "github.com/confyrm/gorest/slack"
  . "github.com/confyrm/gorest/errors"
  "github.com/confyrm/gorest/config"
  . "github.com/confyrm/gorest/servers/slack/commands"
)



// Make a static decoder for performance
var decoder = schema.NewDecoder()
var commandRouter = SlashCommands.New()

// This is set the first time HelpResponse is called.  It is set by Loading
// The help text file located in ../../../help/help.hcl
var helpResponses map[string]interface{}

// SlashRouter is the top level slash command router.
func SlashRouter(config *config.Config, rw http.ResponseWriter, req *http.Request) error {

  if err := req.ParseForm(); err != nil {
      return StatusError{http.StatusBadRequest,
        errors.New(fmt.Sprintf("Could not parse form data: %#v", err.Error()))}
  }

  sReq := &slack.Request{}

  if err := decoder.Decode(sReq, req.PostForm); err != nil {
    return StatusError{http.StatusBadRequest,
      errors.New(fmt.Sprintf("Unrecognized content: %#v", err.Error()))}
  }
  // Dump the slack.Request to the log
  sReq.Log()

  command, err := sReq.TextToCommand()
  if err != nil {
    return StatusError{http.StatusInternalServerError, err}
  }
  log.Printf("Received Slack slash command: %#v", command)

  // First see if this is a help request.  If so, return a help response
  if (len(command.Commands) == 0 && len(command.Params) == 0) ||
    command.Commands[len(command.Commands)-1] == "help" {
    rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
    response := HelpResponse(config, command)
    if err := json.NewEncoder(rw).Encode(response); err != nil {
      return StatusError{http.StatusInternalServerError, err}
    }
    return nil
  }

  // First authZ test.  Ensure the token provided in the slack request
  // matches the SLACK_TOKEN
  if sReq.Token != config.GetString("SLACK_TOKEN") {
    return StatusError{http.StatusUnauthorized,
      errors.New("Not authorized. Wrong Slack Token.")}
  }

  // Get the route handler for this command
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
func HelpResponse(config *config.Config, command *slack.DevHubCommand) *slack.Response {
  // Although called every time, this will only attempted to read the help
  // text the first time its called.
  ImportHelpText(config)

  cLen := len(command.Commands)
  var text string
  // If the only command is help, then return the top level help
  if cLen == 1 && command.Commands[0] == "help" {
    if tryit, ok := helpResponses["base"]; ok {
      text = tryit.(string)
    }
  } else {
    // Create the help lookup key, by joining all the commands except the last
    // one, which should be "help"
    key := strings.Join(command.Commands[:cLen - 1], ".")
    t, ok := helpResponses[key]
    if !ok {
      // No help for this command.  Just return the top level help.
      log.Printf("Help request received for %s, but no help found\n", key)
      if tryit, ok := helpResponses["base"]; ok {
        text = tryit.(string)
      }
    } else {
      text = t.(string)
    }
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

  helpPath := filepath.Join(config.GetString("APP_ROOT"), "help", "help.hcl" )
  log.Printf("Reading help from %s", helpPath)
  var err error
  helpResponses, err = ParseHelpText(helpPath)
  if err != nil {
    log.Printf("Error parsing help file: ", err)
    helpResponses = make(map[string]interface{})
  }
  if helpResponses == nil {
    log.Printf("Error parsing help file. No map returned: ", helpResponses)
    helpResponses = make(map[string]interface{})
  }
}

func ParseHelpText(helpPath string) (map[string]interface{}, error) {

  helpText, err := ioutil.ReadFile(helpPath)
  if err != nil {
    log.Printf("Could not load help text from %s: %#v", helpPath, err)
    return nil, err
  }
  var out interface{}
  err = hcl.Decode(&out, string(helpText))
  if err != nil {
    log.Printf("Could not parse help text in %s: %#v", helpPath, err)
    return nil, err
  }
  return out.(map[string]interface{}), nil
}
