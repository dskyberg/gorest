package cmd

import (
  "fmt"
  "log"
  "net/http"
  "encoding/json"
  "errors"
  "strings"

  "golang.org/x/oauth2"
  "github.com/google/go-github/github"

  . "github.com/confyrm/gorest/errors"
  "github.com/confyrm/gorest/config"
  "github.com/confyrm/gorest/slack"
)

type DevHubCommand struct {
  Commands []string
  Params map[string]string
}

const githubUsage = `
  github accepts the following commands:
  new
  get
  update
  close
`

const newUsage = `
  github new accepts a set of <key>=<value> statements.  Do not use
  the '=' character within a value.  Unrecognized keys are ignored.
  Here are the supported keys:
  -
  title:      Required. Issue title
  body:       Optional. Text to add to the issue body
  labels:     Optionall A comma separated list of labels
  milestone:  Optional. The milestone must exist in the identified repository
  assignee:   Optional.  Github user to assign the issue to
  repository: Optional. If not provided, devhub will be used.
  -
  Example: /github new title = Here is my title labels = label1, label2
`

func DevHub(config *config.Config, sReq *slack.Request, rw http.ResponseWriter, req *http.Request) error {

  rw.Header().Set("Content-Type", "application/json; charset=UTF-8")

  /*
  text := fmt.Sprintf("Hi %s! Your [%s] request is being processed!", sReq.UserName, sReq.Command)
  response := slack.Response{slack.Ephemeral, text, nil}
  if err := json.NewEncoder(rw).Encode(response); err != nil {
      return err
  }
  */

  command, err := TextToCommand(sReq.Text)
  if err != nil {
    return StatusError{http.StatusInternalServerError, err}
  }

  if command.Commands[0] == "help" {
    return HelpResponse(rw, "devhub")
  }

  text, err := AddTicket(config, *command)
  if err != nil {
    return StatusError{http.StatusInternalServerError, err}
  }

  response := slack.Response{slack.Ephemeral, *text, nil}

  if err := json.NewEncoder(rw).Encode(response); err != nil {
    return StatusError{http.StatusInternalServerError, err}
  }
  //rw.WriteHeader(http.StatusOK)
  return nil
}

/*
  Send a help text response, if the user requested it
*/

func HelpResponse(rw http.ResponseWriter, command string) error {

  var text string
  if command == "devhub" {
    text = githubUsage
  }
  response := slack.Response{slack.Ephemeral, text, nil}

  if err := json.NewEncoder(rw).Encode(response); err != nil {
    return StatusError{http.StatusInternalServerError, err}
  }
  return nil
}


/*
  Break the provided text into components.  Format is
  command part part part
  part: key : value

  github.Issue
  {
    Number:4,
    State:"open",
    Title:"this another title of the ticket",
    User:github.User
    {
      Login:"dskyberg-confyrm",
      ID:14827293,
      AvatarURL:"https://avatars.githubusercontent.com/u/14827293?v=3",
      HTMLURL:"https://github.com/dskyberg-confyrm",
      GravatarID:"",
      Type:"User",
      SiteAdmin:false,
      URL:"https://api.github.com/users/dskyberg-confyrm",
      EventsURL:"https://api.github.com/users/dskyberg-confyrm/events{/privacy}",
      FollowingURL:"https://api.github.com/users/dskyberg-confyrm/following{/other_user}",
      FollowersURL:"https://api.github.com/users/dskyberg-confyrm/followers",
      GistsURL:"https://api.github.com/users/dskyberg-confyrm/gists{/gist_id}",
      OrganizationsURL:"https://api.github.com/users/dskyberg-confyrm/orgs",
      ReceivedEventsURL:"https://api.github.com/users/dskyberg-confyrm/received_events",
      ReposURL:"https://api.github.com/users/dskyberg-confyrm/repos",
      StarredURL:"https://api.github.com/users/dskyberg-confyrm/starred{/owner}{/repo}",
      SubscriptionsURL:"https://api.github.com/users/dskyberg-confyrm/subscriptions"
    },
    Labels:
    [
      github.Label
      {
        URL:"https://api.github.com/repos/confyrm/confyrm.github.io/labels/investigating",
        Name:"investigating",
        Color:"1192FC"
      }
      github.Label
      {
        URL:"https://api.github.com/repos/confyrm/confyrm.github.io/labels/EPS",
        Name:"EPS",
        Color:"171717"
      }
    ],
    Comments:0,
    CreatedAt:time.Time
    {
      sec:, nsec:,
      loc:time.Location
      {
        name:"UTC",
        cacheStart:, cacheEnd:
      }
    },
    UpdatedAt:time.Time
    {
      sec:, nsec:,
      loc:time.Location
      {
        name:"UTC",
        cacheStart:, cacheEnd:
      }
    },
    URL:"https://api.github.com/repos/confyrm/confyrm.github.io/issues/4",
    HTMLURL:"https://github.com/confyrm/confyrm.github.io/issues/4"
  }
*/

func AddTicket(config *config.Config, command DevHubCommand) (*string, error) {
  ts := oauth2.StaticTokenSource(
    &oauth2.Token{AccessToken: config.GetString("GITHUB_TOKEN")},
  )
  tc := oauth2.NewClient(oauth2.NoContext, ts)

  client := github.NewClient(tc)

  input := TextToIssueRequest(command)

  issue, _, err := client.Issues.Create("confyrm", "confyrm.github.io", input)
	if err != nil {
    errr := errors.New(fmt.Sprintf("Issue creation failed with %s", err))
    return nil, errr
	}
  log.Printf("Issue created: %+v", *issue)

  text := fmt.Sprintf("Issue %d was successfully created: %s", *issue.Number, *issue.HTMLURL)

  return &text, nil
}

func FormatResponse(config config.Config, sReq slack.Request, command DevHubCommand, issue github.Issue) (*slack.Response, error) {

  // Top level text.  Everything else will be in the attachment
  text := fmt.Sprintf("@%s created an issue", sReq.UserName)


  // If there was no repo listed in the slash command, then use devhub
  var repo string
  if i, ok := command.Params["repo"]; ok {
    repo = i
  } else {
    repo = "devhub"
  }
  log.Printf("Using repo: %s", repo)

  attachment := slack.Attachment{}

  response := slack.Response{ResponseType: slack.InChannel, Text: text, Attachments: slack.Attachments{ attachment} }

  return &response, nil
}

func GetUser(client *github.Client) (*string, error) {

  user, _, err := client.Users.Get("")
  if err != nil {
    errr := errors.New(fmt.Sprintf("client.Users.Get() faled with '%s'\n", err))
    return nil, errr
  }
  return user.Login , nil
}

func TextToIssueRequest(command DevHubCommand) *github.IssueRequest {

  input := &github.IssueRequest{}
  if i, ok := command.Params["title"]; ok {
    input.Title = &i
  }
  if i, ok := command.Params["body"]; ok {
    input.Body = &i
  }
  if i, ok := command.Params["assignee"]; ok {
    input.Body = &i
  }

  if i, ok := command.Params["labels"]; ok {
    s := strings.Split(i, ",")
    if len(s) > 0 {
      for j := 0; j < len(s); j++ {
        s[j] = strings.Trim(s[j], " ")
      }
      input.Labels = &s
    }
  }
  return input
}

func TextToCommand(text string) (*DevHubCommand, error) {
  t := strings.Trim(text, " ")
  if len(t) == 0 {
    return nil, errors.New("Error parsing DevHub command: Empty string")
  }

  commands, kvText := ParseCommands(t)

  //Now split on '='.  This will yield an even number of values in s, where
  // s[i] = key, s[i+1] = value
  kv := make(map[string]string)

  for x := 0; x < 10; x++ { // Just to prevent infinite loop!
    i := strings.Index(kvText, "=")

    if i == -1 {
      break
    }

    k := strings.ToLower(strings.Trim(kvText[:i], " "))
    fmt.Printf(" - k: %s\n", k)
    kvText = strings.Trim(kvText[i + 1:], " ")

    if len(kvText) == 0 {
      return nil, errors.New(fmt.Sprintf("Error parsing KV pairs: No remainder for key: %s ", k ))
    }

    // i = strings.Index(kvText, "=")
    j := FindStartOfNextKey(kvText)
    if j == -1 {
      // Last value in k/v pairs
      v := strings.Trim(kvText, " ")
      kv[k] = v
      fmt.Printf(" - v: [%s]\n", v)
      break
    }


    // j now sits at letter of the next key.
    v := strings.ToLower(strings.Trim(kvText[:j - 1], " "))
    fmt.Printf(" - v: [%s]\n", v)
    kv[k] = v
    kvText = strings.Trim(kvText[j:], " ")
  }

  return &DevHubCommand{commands, kv}, nil
}

func FindStartOfNextKey(text string) int {
  t := strings.TrimLeft(text, " ")

  if len(t) == 0 {
    // No text to test
    return -1
  }
  i := strings.Index(t, "=")

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

  for i--; i > 0; i-- {
    if t[i] != ' ' {
      break
    }
  }

  // Now skip the key
  for {
    if i == 1 {
      break
    }
    if t[i - 1] == ' ' {
      break
    }
    i--
  }
  return i
}

func ParseCommands(text string) ([]string, string) {
  // Grab everything up to the first key
  var cmdText string
  var kvText string
  firstKey := FindStartOfNextKey(text)

  if firstKey > 0 {
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
    tmp[i] = strings.Trim(tmp[i], " ")
  }
  var commands []string
  for _, x := range tmp {
    if len(x) > 0 {
      commands = append(commands, x)
    }
  }
  return commands, kvText
}
