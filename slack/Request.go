package slack

import (
  "log"
  "errors"
  "fmt"
  "net/http"

  "github.com/gorilla/schema"
)
/*
2016/07/07 23:08:41
token: [5XhdiDjSY9C9EwCjVS5LwKpv]
team_id: [T0A06JCCQ]
team_domain: [confyrm]
channel_id: [C1PMNSN2Z]
channel_name: [test]
user_id: [U0A06NUP5]
user_name: [david]
command: [/devhub]
text: [help]
response_url: [https://hooks.slack.com/commands/T0A06JCCQ/57857840896/QTIerV9ZtswZHS36ijsav2Kr]
*/

// The Slack Request that is sent by Slack to a custom app.
type Request struct {
  // The token that was used by Slack.  This should be validated against
  // SLACK_TOKEN.
  Token string `schema:"token"`
  // The Slack team the below user is part of.
  TeamId string `schema:"team_id"`
  // The Slack team domain.
  TeamDomain string `schema:"team_domain"`
  // The Slack channel the slash command was executed on.
  ChannelId string `schema:"channel_id"`
  // The name of the Slack channel the slash command was executed on.
  ChannelName string `schema:"channel_name"`
  // The Id of the user executing the slash command.
  UserId string `schema:"user_id"`
  // The Slack user name of the user that executed the slash command.
  UserName string `schema:"user_name"`
  // The Slack slash command that was executed.
  Command string `schema:"command"`
  // all text that was provided after the slash command.
  Text string `schema:"text"`
  // If the time to respond to the command is long, the response can be
  // sent to this url, instead ofo in the initial response.
  ResponseUrl string `schema:"response_url"`
}

var decoder = schema.NewDecoder()

// Decode does a JSON decode of the provided map.  This is
// Generally passed the url.Values from a http.Request.
func (r *Request) DecodeHttp(req *http.Request) error {
  if req == nil {
    return errors.New("req is nil")
  }
  // Make sure there's a PostForm available
  if req.PostForm == nil {
    if err := req.ParseForm(); err != nil {
      errors.New(fmt.Sprintf("Could not parse form data: %#v", err.Error()))
    }
  }
  // Just pass to
  if err := r.DecodeMap(req.PostForm); err != nil {
    return errors.New(fmt.Sprintf("Unrecognized content: %#src", err.Error()))
  }
  return nil
}

// Decode does a JSON decode of the provided map.  This is
// Generally passed the url.Values from a http.Request.
func (r *Request) DecodeMap(src map[string][]string) error {
  if err := decoder.Decode(r, src); err != nil {
    return errors.New(fmt.Sprintf("Unrecognized content: %#src", err.Error()))
  }
  return nil
}


func (r *Request) Log() {
  log.Printf("token: [%s] team_id: [%s] team_domain: [%s] channel_id: [%s] channel_name: [%s] user_id: [%s] user_name: [%s] command: [%s] text: [%s] response_url: [%s]",
    r.Token, r.TeamId, r.TeamDomain, r.ChannelId, r.ChannelName, r.UserId, r.UserName, r.Command, r.Text, r.ResponseUrl)
}
