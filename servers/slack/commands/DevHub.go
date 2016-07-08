package commands

import (
  "fmt"
  "log"
  "errors"
  "strings"
  "strconv"
  "golang.org/x/oauth2"
  "github.com/google/go-github/github"
  //"gopkg.in/libgit2/git2go.v22"
  . "github.com/confyrm/gorest/errors"
  "github.com/confyrm/gorest/config"
  "github.com/confyrm/gorest/slack"
  "github.com/confyrm/gorest/githubclient"
)

func DevHub(config *config.Config, sReq *slack.Request, command *slack.DevHubCommand) (*slack.Response, *StatusError) {
  log.Printf("DevHub was called")

  if len(command.Commands) < 1 {
    RespondWithError(sReq, errors.New("No devhub command specified"))
  }

  var (
    resp *slack.Response
    err error
  )

  cmd := command.Commands[0]
  switch (cmd) {
    case "new":
      resp, err = HandleNew(sReq, config, command)
    case "get":
      resp, err = HandleGet(sReq, config, command)
    default:
      err = fmt.Errorf("Command not recognized: %s", cmd)
  }

  if err != nil {
    RespondWithError(sReq, err)
  } else {
    RespondWithSuccess(sReq, resp)
  }

  // Since this is a long running command, just return (nil, nil)
  return nil, nil
}

func RespondWithSuccess(sReq *slack.Request, response *slack.Response ) {

  if err := sReq.Respond(response); err != nil {
    log.Printf("Error sending error response: %s\n\n", err.Error())
  }
}

func RespondWithError(sReq *slack.Request, err error) {
  var atts = slack.Attachments {
    slack.Attachment {
      Title: "Oh snap! Something went wrong!",
      Text: err.Error(),
      Color: slack.DANGER,
    },
  }
  response := slack.Response{slack.Ephemeral.String(), nil, atts}
  if errr := sReq.Respond(&response); err != nil {
    log.Printf("Error sending error response: %s %s", err.Error(), errr.Error())
  }
}

func HandleGet(sReq *slack.Request, config *config.Config, command *slack.DevHubCommand) (*slack.Response, error) {

  // Look in the command for the owner and repo.  If none are provided, and no
  // default was provided in the config, then it's an error
  owner := config.GetString(githubclient.DefaultOwner)
  if len(owner) == 0 {
    return nil, errors.New("Could not find a GitHub owner in the command or the config")
  }
  repo := command.ValueOrDefault("repo", config.GetString(githubclient.DefaultRepo))
  if len(repo) == 0 {
    return nil, errors.New("Could not find a GitHub repo in the command or the config")
  }

  // Validate that we have all required data
  var numberStr string
  if value, ok := command.Value("number"); ok {
    numberStr = value
  } else {
    remaining := command.CommandsFrom(1)
    if len(remaining) == 0 {
      // Number was not provided as a KV value, or in the
      // list of Commands.
      return nil, fmt.Errorf("Issue number was not provided")
    }
    // Make sure it's an int
    numberStr = remaining[0]
  }
  number, err := strconv.Atoi(numberStr)
  if err != nil {
    // Couldn't parse the provided number
    return nil, fmt.Errorf("Value provided for number was not an int: %s", numberStr)
  }

  client := NewGithubClient(config)

  issue, _, err := client.Issues.Get(owner, repo, number)
	if err != nil {
    errr := fmt.Errorf("Issue fetch failed with %s", err.Error())
    return nil, errr
	}
  log.Printf("Issue fetched: %+v\n\n", *issue)

  var body string
  if issue.Body == nil || len(*issue.Body) > 0 {
    body = *issue.Body
  } else {
    body = ""
  }
  var assignee string

  if issue.Assignee == nil {
    assignee = ""
  } else {
    assignee = fmt.Sprintf("<%s|%s>", issue.Assignee.HTMLURL, issue.Assignee.Name)
  }
  createdBy := fmt.Sprintf("<%s|%s>", issue.User.HTMLURL, issue.User.Name)
  var milestone string
  if issue.Milestone == nil {
    milestone = ""
  } else {
    milestone = fmt.Sprintf("%d",issue.Milestone.Number)
  }
  var atts = slack.Attachments {
    slack.Attachment {
      Fallback: fmt.Sprintf("#%d: %s\n%s", *issue.Number, *issue.HTMLURL, *issue.Title),
      Title: fmt.Sprintf("<%s|#%d>: %s",*issue.HTMLURL, *issue.Number, *issue.Title),
      Text: body,
      Color: slack.GOOD,
      MarkdownIn: []string{"title", "text"},
    },
    slack.Attachment {
      Title: "Details",
      Text: fmt.Sprintf("- Status: %s\n- Created by %s\n- Assigned: %s\n- Milestone: %s",
        issue.State, createdBy, assignee, milestone),
      MarkdownIn: []string{"title", "text"},
    },
  }
  title := "Get Issue"
  response := slack.Response{slack.Ephemeral.String(), &title, atts}
  return &response, nil
}

// HandleNew creates a new GitHub issue, using the Key/Value data in the DevHubCommand
func HandleNew(sReq *slack.Request, config *config.Config, command *slack.DevHubCommand) (*slack.Response, error) {

  // Look in the command for the owner and repo.  If none are provided, and no
  // default was provided in the config, then it's an error
  owner := config.GetString(githubclient.DefaultOwner)
  if len(owner) == 0 {
    return nil, errors.New("Could not find a GitHub owner in the command or the config")
  }
  repo := command.ValueOrDefault("repo", config.GetString(githubclient.DefaultRepo))
  if len(repo) == 0 {
    return nil, errors.New("Could not find a GitHub repo in the command or the config")
  }

  // Validate that we have all required data
  if !command.HasValue("title") {
    return nil, fmt.Errorf("Issue title was not provided")
  }

  client := NewGithubClient(config)
  input := TextToIssueRequest(command)

  issue, _, err := client.Issues.Create(owner, repo, input)
	if err != nil {
    errr := fmt.Errorf("Issue creation failed with %s", err.Error())
    return nil, errr
	}
  log.Printf("Issue created: %+v", *issue)

  var body string
  if issue.Body == nil || len(*issue.Body) > 0 {
    body = *issue.Body
  } else {
    body = ""
  }
  var atts = slack.Attachments {
    slack.Attachment {
      Title: fmt.Sprintf("#%d: %s", *issue.Number, *issue.Title),
      TitleLink: *issue.HTMLURL,
      Text: body,
      Color: slack.GOOD,
    },
    slack.Attachment {
      Title: fmt.Sprintf("Created by <@%s|%s>",  issue.User.Name),
    },
  }
  title := fmt.Sprintf("<@%s|%s> created a new issue!", sReq.UserId, sReq.UserName)
  response := slack.Response{slack.InChannel.String(), &title, atts}
  return &response, nil
}

// NewGithubClient is a utility function that uses the GITHUB_TOKEN from
// the config to create a GitHub client.
func NewGithubClient(config *config.Config) *github.Client {
  ts := oauth2.StaticTokenSource(
    &oauth2.Token{AccessToken: config.GetString("GITHUB_TOKEN")},
  )
  tc := oauth2.NewClient(oauth2.NoContext, ts)
  client := github.NewClient(tc)
  return client
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
  if i, ok := command.Value("title"); ok {
    input.Title = &i
  }
  if i, ok := command.Value("body"); ok {
    input.Body = &i
  }
  if i, ok := command.Value("assignee"); ok {
    input.Body = &i
  }

  if i, ok := command.Value("labels"); ok {
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

/*
  Returned from Issue.Create
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

  Returned from Issue.Get

  github.Issue{
    Number:617,
    State:"open",
    Title:"Don't run ITs which generate random entities each time, but run them nightly instead",
    Body:"We have 2 ITs marked with "groups={ generated_entities }" annotation. We need to skip these tests by group in Jenkins in all builds except for nightly tests run, so we:
    - only generate 2 sets of random entities on DEV per day (less load on Kafka)
    - speedup usual ITs builds, but keep testing with generated entities also",
    User:github.User{
      Login:"daniilyar-confyrm",
      ID:14834437,
      AvatarURL:"https://avatars.githubusercontent.com/u/14834437?v=3",
      HTMLURL:"https://github.com/daniilyar-confyrm",
      GravatarID:"", T
      ype:"User",
      SiteAdmin:false,
      URL:"https://api.github.com/users/daniilyar-confyrm",
      EventsURL:"https://api.github.com/users/daniilyar-confyrm/events{/privacy}",
      FollowingURL:"https://api.github.com/users/daniilyar-confyrm/following{/other_user}",
      FollowersURL:"https://api.github.com/users/daniilyar-confyrm/followers",
      GistsURL:"https://api.github.com/users/daniilyar-confyrm/gists{/gist_id}",
      OrganizationsURL:"https://api.github.com/users/daniilyar-confyrm/orgs",
      ReceivedEventsURL:"https://api.github.com/users/daniilyar-confyrm/received_events",
      ReposURL:"https://api.github.com/users/daniilyar-confyrm/repos",
      StarredURL:"https://api.github.com/users/daniilyar-confyrm/starred{/owner}{/repo}",
      SubscriptionsURL:"https://api.github.com/users/daniilyar-confyrm/subscriptions"
    },
    Labels:[
      github.Label{
        URL:"https://api.github.com/repos/confyrm/devhub/labels/1sp",
        Name:"1sp",
        Color:"FBDDCE"
      }
      github.Label{
        URL:"https://api.github.com/repos/confyrm/devhub/labels/Todo",
        Name:"Todo",
        Color:"ededed"
      }
    ],
    Assignee:github.User{
      Login:"daniilyar-confyrm",
      ID:14834437,
      AvatarURL:"https://avatars.githubusercontent.com/u/14834437?v=3",
      HTMLURL:"https://github.com/daniilyar-confyrm",
      GravatarID:"",
      Type:"User",
      SiteAdmin:false,
      URL:"https://api.github.com/users/daniilyar-confyrm",
      EventsURL:"https://api.github.com/users/daniilyar-confyrm/events{/privacy}",
      FollowingURL:"https://api.github.com/users/daniilyar-confyrm/following{/other_user}",
      FollowersURL:"https://api.github.com/users/daniilyar-confyrm/followers",
      GistsURL:"https://api.github.com/users/daniilyar-confyrm/gists{/gist_id}",
      OrganizationsURL:"https://api.github.com/users/daniilyar-confyrm/orgs",
      ReceivedEventsURL:"https://api.github.com/users/daniilyar-confyrm/received_events",
      ReposURL:"https://api.github.com/users/daniilyar-confyrm/repos",
      StarredURL:"https://api.github.com/users/daniilyar-confyrm/starred{/owner}{/repo}",
      SubscriptionsURL:"https://api.github.com/users/daniilyar-confyrm/subscriptions"
    },
    Comments:0,
    CreatedAt:time.Time{sec:, nsec:, loc:time.Location{name:"UTC", cacheStart:, cacheEnd:}},
    UpdatedAt:time.Time{sec:, nsec:, loc:time.Location{name:"UTC", cacheStart:, cacheEnd:}},
    URL:"https://api.github.com/repos/confyrm/devhub/issues/617",
    HTMLURL:"https://github.com/confyrm/devhub/issues/617",
    Milestone:github.Milestone{
      URL:"https://api.github.com/repos/confyrm/devhub/milestones/11",
      Number:11,
      State:"open",
      Title:"sprint_30",
      Description:"2016-04-27 - 2016-05-10",
      Creator:github.User{
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
      OpenIssues:9,
      ClosedIssues:35,
      CreatedAt:time.Time{sec:, nsec:, loc:time.Location{name:"UTC", cacheStart:, cacheEnd:}},
      UpdatedAt:time.Time{sec:, nsec:, loc:time.Location{name:"UTC", cacheStart:, cacheEnd:}}
    }
  }
*/
