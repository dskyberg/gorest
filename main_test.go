package main

import (
  "testing"
  "net/url"
  "net/http"
  "io/ioutil"
  "net/http/httptest"

  "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/require"
  "github.com/confyrm/gorest/router/handler"
  "github.com/confyrm/gorest/servers/slack/routes"
)

var form = url.Values {
  "token": {"abcd1234"},
  "team_domain": {"example"},
  "channel_id": {"C2147483705"},
  "channel_name": {"test"},
  "user_id": {"U2147483697"},
  "user_name": {"{Steve"},
  "command": {"/devhub"},
  "text": {"help"},
  "response_url": {"http://localhost:8002/test"},
}


func TestMain(t *testing.T ) {

  c := SetupConfig()
  h := handler.Handler {
    c,
    routes.SlashRouter,
  }

  server := httptest.NewServer(h)
  defer server.Close()

  resp, err := http.PostForm(server.URL, form)
  require.Nil(t, err, "PostForm returned nil")
  assert.Equal(t, 200, resp.StatusCode, "Bad response %d", resp.StatusCode)
  body, err := ioutil.ReadAll(resp.Body)
  resp.Body.Close()
  t.Logf("Response body: %s", body)
}
