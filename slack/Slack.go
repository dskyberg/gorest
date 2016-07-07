// Package slack provides a set of types for using the Slack API. This includes
// types for custom commands to receive REST requests from Slack.
package slack

import (
  "fmt"
  "log"
  "strings"
  "errors"
)
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

func (r *Request) Log() {
  log.Printf("token: [%s] team_id: [%s] team_domain: [%s] channel_id: [%s] channel_name: [%s] user_id: [%s] user_name: [%s] command: [%s] text: [%s] response_url: [%s]",
    r.Token, r.TeamId, r.TeamDomain, r.ChannelId, r.ChannelName, r.UserId, r.UserName, r.Command, r.Text, r.ResponseUrl)
}

// InChannel tells Slack the response should be placed in the channel for all
// channel members to see.
const InChannel = "in_channel"
// Ephemeral tells Slack the response should be displayed only to the user that
// executed the slash command. For long running commands, consider immediately
// responding to the user with an Ephemeral response, and use InChannel for the
// final result.
const Ephemeral = "ephemeral"

// Slash standard green color for attachments.
const GOOD = "good"
// Slash standard yellow color for attachments.
const WARNING = "warning"
// Slash standard red color for attachments.
const DANGER = "danger"

// AttachmentFields are displayed as tabular data.
type AttachmentField struct {
  // Table header value.
  Title string `json:"title"`
  // Table value.  Displayed under the header.
  Value string `json:"value"`
  // true, if this field is short enough to be displayed on same line as
  // previous fields.  Use false, if you know it is too long
  Short bool `json:"short"`
}
// Helper type for a slice of AttachmentFields.
type AttachmentFields []AttachmentField

// Attachment is the primary tool for creating rich responses.
type Attachment struct {
  Title string `json:"title"`
  TitleLink string `json:"title_link, omitempty"`

  PreText string `json:"pretext, omitempty"`
  Text string `json:"text, omitempty"`
  Color string `json:"color, omitempty"`
  Fallback string `json:"fallback"`

  ImageUrl string `json:"image_url, omitempty"`
  ThumbUrl string `json:"thumb_url, omitempty"`

  // Author fields
  AuthorName string `json:"author_name, omitempty"`
  AuthorLink string `json:"author_link, omitempty"`
  AuthorIcon string `json:"author_icon, omitempty"`


  // Footer fields
  Footer string `json:"footer, omitempty"`
  FooterIcon string `json:"footer_icon, omitempty"`
  TimeStamp int `json:"ts", omitempty`

  // Table fields
  Fields AttachmentFields `json:"fields, omitempty"`
}

// Helper type for a slick of Attachments.
type Attachments []Attachment

// The Response must be json encoded.  Response data must be URL encoded, also.
type Response struct {
  ResponseType string `json:"response_type"`
  Text string `json:"text"`
  Attachments Attachments `json:"attachments,omitempty"`
}

// The Slack slash commands that we process contain a consistent structure that
// can be leveraged for any type of activity.
type DevHubCommand struct {
  Commands []string
  Params map[string]string
}

const SLASH_DELIM = " \n"
const TRIM_CUTSET = " \n"
const KV_DELIM = "="


func (sReq *Request) TextToCommand() (*DevHubCommand, error) {
  // Trim the string first, to remove any unwanted spaces and new lines
  t := strings.Trim(sReq.Text, TRIM_CUTSET)

  // Get the set of commands, and whatever text may be remaining after the commands
  commands, kvText := ParseCommands(t)
  kv, err := ParseKeyValuePairs(kvText)
  if err != nil {
    return nil, err
  }

  return &DevHubCommand{commands, kv}, nil
}

// ParseCommands is a helper function that parses out any commands that are
// placed before the Key/Value pairs in the provided text.
// Note, if there are no commands AND no KV pairs, ParseCommands will return
// an empty string for kvText.
func ParseCommands(text string) ([]string, string) {
  // Grab everything up to the first key
  var cmdText string
  var kvText string
  commands := make([]string, 0, 10)
  firstKey := FindStartOfNextKey(text)

  if firstKey == 0 {
    // There are KV's but no commands
    return commands, text
  } else if firstKey > 0 {
    cmdText = text[:firstKey-1]
    kvText = text[firstKey:]
  } else {
    // There doesn't appear to be any kv pairs.
    cmdText = text
    kvText = ""
  }
  // Read the list of commands first
  tmp := strings.Split(cmdText, " ")
  for i := 0; i < len(tmp); i++ {
    tmp[i] = strings.Trim(tmp[i], TRIM_CUTSET)
  }
  for _, x := range tmp {
    if len(x) > 0 {
      commands = append(commands, x)
    }
  }
  return commands, kvText
}

func ParseKeyValuePairs(kvText string) (map[string]string, error) {
  //Split on KV_DELIM.  This will yield an even number of values in s, where
  // s[i] = key, s[i+1] = value
  kv := make(map[string]string)

  for x := 0; x < 10; x++ { // Just to prevent infinite loop!
    i := strings.Index(kvText, KV_DELIM)

    if i == -1 {
      break
    }

    k := strings.ToLower(strings.Trim(kvText[:i], TRIM_CUTSET))
    //fmt.Printf(" - k: %s\n", k)
    kvText = strings.Trim(kvText[i + 1:], TRIM_CUTSET)

    if len(kvText) == 0 {
      err := errors.New(fmt.Sprintf("Error parsing Key Value pairs: No remainder for key: %s ", k ))
      return kv, err
    }

    j := FindStartOfNextKey(kvText)
    if j == -1 {
      // Last value in k/v pairs
      v := strings.Trim(kvText, TRIM_CUTSET)
      kv[k] = v
      //fmt.Printf(" - v: [%s]\n", v)
      break
    }


    // j now sits at letter of the next key.
    v := strings.ToLower(strings.Trim(kvText[:j - 1], TRIM_CUTSET))
    //fmt.Printf(" - v: [%s]\n", v)
    kv[k] = v
    kvText = strings.Trim(kvText[j:], TRIM_CUTSET)
  }
  return kv, nil
}


func FindStartOfNextKey(text string) int {
  t := strings.TrimLeft(text, SLASH_DELIM)

  if len(t) == 0 {
    // No text to test
    return -1
  }
  i := strings.Index(t, KV_DELIM)

  if i == -1 {
    // No keys in the text
    return i
  }

  if i == 0 {
    // No keys, because text starts with '='
    return -1
  }

  if i == len(t) - 1 {
    // No keys, because text ends with '='
    return -1
  }

  // skip over any white space
  for i--; i > 0; i-- {
    if t[i] != ' ' && t[i] != '\n' && t[i] != '\r'  && t[i] != '\f' {
      break
    }
  }

  // Now skip the key
  for {
    if i == 0 {
      break
    }
    if t[i - 1] == ' ' {
      break
    }
    i--
  }
  return i
}
