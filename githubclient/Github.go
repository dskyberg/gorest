// Package github provides an interface for interacting with GitHub.
package githubclient

import (
  "fmt"
  "errors"
  "golang.org/x/oauth2"
  "github.com/google/go-github/github"
)

const DefaultOwner = "GITHUB_DEFAULT_OWNER"
const DefaultRepo = "GITHUB_DEFAULT_REPO"

type GithubClient struct {
    ClientId string
    ClientSecret string
    Token string
    // When a GITHUB_TOKEN is available, we can use a static token source
    ts *oauth2.TokenSource
    client *github.Client
}

func New(token string) *GithubClient {

  ts := oauth2.StaticTokenSource(
    &oauth2.Token{AccessToken: token},
  )
  tc := oauth2.NewClient(oauth2.NoContext, ts)

  client := github.NewClient(tc)

  var gh = GithubClient{
    "",
    "",
    token,
    &ts,
    client,
  }
  return &gh
}

func (client GithubClient) AddTicket(owner string, repo string, ir *github.IssueRequest) (*github.Issue, error) {

  issue, _, err := client.client.Issues.Create(owner, repo, ir)
	if err != nil {
    errr := errors.New(fmt.Sprintf("Issue creation failed with %s", err))
    return nil, errr
	}

  return issue, nil
}
