package slack

import (
  "log"
  "bytes"
  "io/ioutil"
  "net/http"
  "crypto/tls"
  "encoding/json"
  "errors"
  "fmt"
)
// Respond is used by long running commands to send a Response to the
// URL provided in the slack.Request It's just a wrapper around
// http.Request
func (sReq *Request) Respond(sResp *Response) error {

  if sReq.ResponseUrl == "" {
    return errors.New("No ResponseUrl in Request")
  }

  buffer := new(bytes.Buffer)
  json.NewEncoder(buffer).Encode(*sResp)
  tr := &http.Transport{
      TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
      DisableKeepAlives: true,
  }
  client := &http.Client{Transport: tr}
  resp, err := client.Post(sReq.ResponseUrl,
    "application/json; charset=utf-8", buffer)

  if err != nil {
    // Error sending the post
    return fmt.Errorf("Error posting response: %s", err.Error())
  }

  defer resp.Body.Close()
  //block forever at the next line
  content, _ := ioutil.ReadAll(resp.Body)

  log.Printf("Response %d: %s %s\n\n",resp.StatusCode, buffer.String(), string(content))

  return nil
}
