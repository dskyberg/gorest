package cmd

import (
  "fmt"
  "log"
  "net/http"
  "errors"
  "strings"

  "golang.org/x/oauth2"
  "github.com/google/go-github/github"

  . "github.com/confyrm/gorest/errors"
  "github.com/confyrm/gorest/config"
  "github.com/confyrm/gorest/slack"
)

func DevHub(config *config.Config, sReq *slack.Request) (*slack.Response, *StatusError) {

  command, err := sReq.TextToCommand()
  if err != nil {
    return nil, &StatusError{http.StatusInternalServerError, err}
  }

  text, err := AddTicket(config, command)
  if err != nil {
    return nil, &StatusError{http.StatusInternalServerError, err}
  }

  response := slack.Response{slack.Ephemeral, *text, nil}
  /*
  if err := json.NewEncoder(rw).Encode(response); err != nil {
    return StatusError{http.StatusInternalServerError, err}
  }
  //rw.WriteHeader(http.StatusOK)
  */
  return &response, nil
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

func AddTicket(config *config.Config, command *slack.DevHubCommand) (*string, error) {
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

func FormatResponse(config config.Config, sReq slack.Request, command *slack.DevHubCommand, issue github.Issue) (*slack.Response, error) {

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

func TextToIssueRequest(command *slack.DevHubCommand) *github.IssueRequest {

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
        s[j] = strings.Trim(s[j], " \n")
      }
      input.Labels = &s
    }
  }
  return input
}
