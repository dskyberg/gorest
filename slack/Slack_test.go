package slack

import (
  "testing"
  "reflect"
)


func TestParseKeyValuePairs(t *testing.T) {

    kv, err := ParseKeyValuePairs("title=test number 7 \nlabels=EPS, otherLabel")
    if err != nil {
      t.Errorf("Got an error: %#v", err)
    }
    if len(kv) == 0 {
      t.Error("ParseKeyValuePairs returned an empty map")
    }
    t.Logf("ParseKeyValue: %#v", kv)
}


func TestBasicCommand(t *testing.T) {
  sReq := Request {
    Token: "abcd1234",
    TeamDomain: "example",
    ChannelId: "C2147483705",
    ChannelName: "test",
    UserId: "U2147483697",
    UserName: "Steve",
    Command: "/devhub",
    Text: "cmd1 cmd2 title=test number 7 labels=EPS, otherLabel",
    ResponseUrl: "http://localhost:8002/test",
  }

  command, err := sReq.TextToCommand()
  t.Logf("Request.TextToCommand: %#v", command)
  if err != nil {
    t.Errorf("Epic fail!! %v", err)
    return
  }

  theCmd := DevHubCommand {
    []string{"cmd1", "cmd2"},
    map[string]string{
      "title": "test number 7",
      "labels": "EPS, otherLabel",
    },
  }
  t.Logf("DevHubCommand: %#v", theCmd)

  if !reflect.DeepEqual(*command, theCmd) {
    t.Error("DeepEqual failed.")
  }
}
/*
func TestNoKV(t *testing.T) {
  sReq := Request {
    Token: "abcd1234",
    TeamDomain: "example",
    ChannelId: "C2147483705",
    ChannelName: "test",
    UserId: "U2147483697",
    UserName: "Steve",
    Command: "/devhub",
    Text: "cmd1    cmd2  ",
    ResponseUrl: "http://localhost:8002/test",
  }

  theCmd := DevHubCommand {
    []string{"cmd1","cmd2"},
    map[string]string{
      "title": "test number 7",
      "labels": "EPS, otherLabel",
    },
  }

  command, err := sReq.TextToCommand()
  if err != nil {
    t.Errorf("Epic fail!! %v", err )
    return
  }
  if !reflect.DeepEqual(command, theCmd) {
    t.Error("DeepEqual failed")
  }
}
*/
