package slack

import (
  "testing"
  "reflect"
)


func TestParseCommands_1Command_0KV(t *testing.T) {
  t.Log("ParseCommands: 1 command, 0 KVs")
  startingText := "cmd1"
  goodCommands := []string {"cmd1"}
  goodKvText := ""
  commands, kvText := ParseCommands(startingText)
  if !reflect.DeepEqual(commands, goodCommands) {
    t.Errorf("ParseCommands returned the wrong number of commands: %#v", commands)
  }
  if kvText != goodKvText {
    t.Errorf("ParseCommands returned the wrong remainder: %s", kvText)
  }
}

func TestParseCommands_1Command_0KV_White(t *testing.T) {
  t.Log("ParseCommands: 1 command, 0 KVs with white space")
  startingText := "  \ncmd1 \n  \n  "
  goodCommands := []string {"cmd1"}
  goodKvText := ""
  commands, kvText := ParseCommands(startingText)
  if !reflect.DeepEqual(commands, goodCommands) {
    t.Errorf("ParseCommands returned the wrong number of commands: %#v", commands)
  }
  if kvText != goodKvText {
    t.Errorf("ParseCommands returned the wrong remainder: %s", kvText)
  }

}

func TestParseCommands_2Commands_0KV(t *testing.T) {
  t.Log("ParseCommands: 2 commands, 0 KV")
  startingText := "cmd1 cmd2"
  goodCommands := []string {"cmd1", "cmd2"}
  goodKvText := ""
  commands, kvText := ParseCommands(startingText)
  if !reflect.DeepEqual(commands, goodCommands) {
    t.Errorf("ParseCommands returned the wrong number of commands: %#v", commands)
  }
  if kvText != goodKvText {
    t.Errorf("ParseCommands returned the wrong remainder: %s", kvText)
  }
}

func TestParseCommands_2Commands_0KV_White(t *testing.T) {
  t.Log("ParseCommands: 2 commands, 0 KV")
  startingText := "  \ncmd1 \n  \n    \ncmd2 \n  \n  "
  goodCommands := []string {"cmd1", "cmd2"}
  goodKvText := ""
  commands, kvText := ParseCommands(startingText)
  if !reflect.DeepEqual(commands, goodCommands) {
    t.Errorf("ParseCommands returned the wrong number of commands: %#v", commands)
  }
  if kvText != goodKvText {
    t.Errorf("ParseCommands returned the wrong remainder: %s", kvText)
  }
}

func TestParseCommands_0Commands_1KV(t *testing.T) {
  t.Log("ParseCommands: 0 commands, 1 KV ")
  startingText := "title=test number 7"
  goodCommands := []string(nil)
  goodKvText := "title=test number 7"
  commands, kvText := ParseCommands(startingText)
  if !reflect.DeepEqual(commands, goodCommands) {
    t.Errorf("ParseCommands returned the wrong number of commands: %#v", commands)
  }
  if kvText != goodKvText {
    t.Errorf("ParseCommands returned the wrong remainder: %s", kvText)
  }
}

func TestParseCommands_0Commands_1KV_White(t *testing.T) {
  t.Log("ParseCommands: 0 commands, 1 KV with white space")
  startingText := "  \n\n \ntitle   \n = test number 7"
  goodCommands := []string(nil)
  goodKvText := "  \n\n \ntitle   \n = test number 7"
  commands, kvText := ParseCommands(startingText)
  if !reflect.DeepEqual(commands, goodCommands) {
    t.Errorf("ParseCommands returned the wrong number of commands: %#v", commands)
  }
  if kvText != goodKvText {
    t.Errorf("ParseCommands returned the wrong remainder: [%s]", kvText)
  }
}

func TestParseKeyValuePairs_0Commands_1KV(t *testing.T) {
  t.Log("ParseKeyValuePairs: 0 commands, 1 KV")
  startingText := "title=test number 7"
  goodKv := map[string]string{
    "title": "test number 7",
  }
  kv, err := ParseKeyValuePairs(startingText)
  if err != nil {
    t.Errorf("ParseKeyValuePairs returned an error: %#v", err)
  }
  if !reflect.DeepEqual(kv, goodKv) {
    t.Errorf("ParseKeyValuePairs returned a bad set of kvs: %#v", kv)
  }
}

func TestParseKeyValuePairs_0Commands_1KV_White(t *testing.T) {
  t.Log("ParseKeyValuePairs: 0 commands, 1 KV")
  startingText := "  \n\n \n  title  \n \n\n =  \n \n\n  \n test number 7  \n \n\n  \n "
  goodKv := map[string]string{
    "title": "test number 7",
  }
  kv, err := ParseKeyValuePairs(startingText)
  if err != nil {
    t.Errorf("ParseKeyValuePairs returned an error: %#v", err)
  }
  if !reflect.DeepEqual(kv, goodKv) {
    t.Errorf("ParseKeyValuePairs returned a bad set of kvs: %#v", kv)
  }
}

func TestParseKeyValuePairs_0Commands_2KV(t *testing.T) {
  t.Log("ParseKeyValuePairs: 0 commands, 1 KV")
  startingText := "title=test number 7 labels=EPS, otherLabel"
  goodKv := map[string]string{
    "title": "test number 7",
    "labels": "EPS, otherLabel",
  }
  kv, err := ParseKeyValuePairs(startingText)
  if err != nil {
    t.Errorf("ParseKeyValuePairs returned an error: %#v", err)
  }
  if !reflect.DeepEqual(kv, goodKv) {
    t.Errorf("ParseKeyValuePairs returned a bad set of kvs: %#v", kv)
  }
}


func TestParseKeyValuePairs_2KV_White(t *testing.T) {
  t.Log("ParseKeyValuePairs: 0 commands, 2 KV with white space")
  startingText := "  \n\n \n  title  \n \n\n =  \n \n\n  \n test number 7  \n \n\n  \n labels  \n \n\n  \n =  \n \n\n  \n EPS, otherLabel  \n \n\n  \n "
  goodKv := map[string]string{
    "title": "test number 7",
    "labels": "EPS, otherLabel",
  }
  kv, err := ParseKeyValuePairs(startingText)
  if err != nil {
    t.Errorf("ParseKeyValuePairs returned an error: %#v", err)
  }
  if !reflect.DeepEqual(kv, goodKv) {
    t.Errorf("ParseKeyValuePairs returned a bad set of kvs: %#v", kv)
  }
}

func TestTextToCommand_1Commands_0_KV(t *testing.T) {
  t.Log("TextToCommand: 1 commands, 0 KV")
  sReq := Request {
    Token: "abcd1234",
    TeamDomain: "example",
    ChannelId: "C2147483705",
    ChannelName: "test",
    UserId: "U2147483697",
    UserName: "Steve",
    Command: "/devhub",
    Text: "cmd1",
    ResponseUrl: "http://localhost:8002/test",
  }

  command, err := sReq.TextToCommand()
  t.Logf("Request.TextToCommand: %#v", command)
  if err != nil {
    t.Errorf("Epic fail!! %v", err)
    return
  }

  theCmd := DevHubCommand {
    []string{"cmd1"},
    map[string]string{},
  }
  t.Logf("DevHubCommand: %#v", theCmd)

  if !reflect.DeepEqual(*command, theCmd) {
    t.Error("DeepEqual failed.")
  }
}

func TestTextToCommand_1Commands_0_KV_White(t *testing.T) {
  t.Log("TextToCommand: 1 commands, 0 KV with white space")
  sReq := Request {
    Token: "abcd1234",
    TeamDomain: "example",
    ChannelId: "C2147483705",
    ChannelName: "test",
    UserId: "U2147483697",
    UserName: "Steve",
    Command: "/devhub",
    Text: "  \n \n\n cmd1  \n \n\n ",
    ResponseUrl: "http://localhost:8002/test",
  }

  command, err := sReq.TextToCommand()
  t.Logf("Request.TextToCommand: %#v", command)
  if err != nil {
    t.Errorf("Epic fail!! %v", err)
    return
  }

  theCmd := DevHubCommand {
    []string{"cmd1"},
    map[string]string{},
  }
  t.Logf("DevHubCommand: %#v", theCmd)

  if !reflect.DeepEqual(*command, theCmd) {
    t.Error("DeepEqual failed.")
  }
}

func TestTextToCommand_2Commands_0_KV(t *testing.T) {
  t.Log("TextToCommand: 2 commands, 0 KV")
  sReq := Request {
    Token: "abcd1234",
    TeamDomain: "example",
    ChannelId: "C2147483705",
    ChannelName: "test",
    UserId: "U2147483697",
    UserName: "Steve",
    Command: "/devhub",
    Text: "cmd1 cmd2",
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
    map[string]string{},
  }
  t.Logf("DevHubCommand: %#v", theCmd)

  if !reflect.DeepEqual(*command, theCmd) {
    t.Error("DeepEqual failed.")
  }
}

func TestTextToCommand_2Commands_0_KV_White(t *testing.T) {
  t.Log("TextToCommand: 2 commands, 0 KV with white space")
  sReq := Request {
    Token: "abcd1234",
    TeamDomain: "example",
    ChannelId: "C2147483705",
    ChannelName: "test",
    UserId: "U2147483697",
    UserName: "Steve",
    Command: "/devhub",
    Text: "  \n \n\n cmd1  \n \n\n   \n \n\n cmd2  \n \n\n ",
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
    map[string]string{},
  }
  t.Logf("DevHubCommand: %#v", theCmd)

  if !reflect.DeepEqual(*command, theCmd) {
    t.Error("DeepEqual failed.")
  }
}

func TestTextToCommand_0Commands_1_KV(t *testing.T) {
  t.Log("TextToCommand: 0 commands, 1 KV")
  sReq := Request {
    Token: "abcd1234",
    TeamDomain: "example",
    ChannelId: "C2147483705",
    ChannelName: "test",
    UserId: "U2147483697",
    UserName: "Steve",
    Command: "/devhub",
    Text: "title=test number 7",
    ResponseUrl: "http://localhost:8002/test",
  }

  command, err := sReq.TextToCommand()
  t.Logf("Request.TextToCommand: %#v", command)
  if err != nil {
    t.Errorf("Epic fail!! %v", err)
    return
  }

  theCmd := DevHubCommand {
    []string(nil),
    map[string]string{
      "title": "test number 7",
    },
  }
  t.Logf("DevHubCommand: %#v", theCmd)

  if !reflect.DeepEqual(*command, theCmd) {
    t.Error("DeepEqual failed.")
  }
}

func TestTextToCommand_0Commands_1_KV_White(t *testing.T) {
  t.Log("TextToCommand: 0 commands, 1 KV with white space")
  sReq := Request {
    Token: "abcd1234",
    TeamDomain: "example",
    ChannelId: "C2147483705",
    ChannelName: "test",
    UserId: "U2147483697",
    UserName: "Steve",
    Command: "/devhub",
    Text: " \n  \n\n title \n  \n\n = \n  \n\n test number 7 \n  \n\n ",
    ResponseUrl: "http://localhost:8002/test",
  }

  command, err := sReq.TextToCommand()
  t.Logf("Request.TextToCommand: %#v", command)
  if err != nil {
    t.Errorf("Epic fail!! %v", err)
    return
  }

  theCmd := DevHubCommand {
    []string(nil),
    map[string]string{
      "title": "test number 7",
    },
  }
  t.Logf("DevHubCommand: %#v", theCmd)

  if !reflect.DeepEqual(*command, theCmd) {
    t.Error("DeepEqual failed.")
  }
}

func TestTextToCommand_0Commands_2_KV(t *testing.T) {
  t.Log("TextToCommand: 0 commands, 2 KV")
  sReq := Request {
    Token: "abcd1234",
    TeamDomain: "example",
    ChannelId: "C2147483705",
    ChannelName: "test",
    UserId: "U2147483697",
    UserName: "Steve",
    Command: "/devhub",
    Text: "title=test number 7 labels=EPS, somethingElse",
    ResponseUrl: "http://localhost:8002/test",
  }

  command, err := sReq.TextToCommand()
  t.Logf("Request.TextToCommand: %#v", command)
  if err != nil {
    t.Errorf("Epic fail!! %v", err)
    return
  }

  theCmd := DevHubCommand {
    []string(nil),
    map[string]string{
      "title": "test number 7",
      "labels": "EPS, somethingElse",
    },
  }
  t.Logf("DevHubCommand: %#v", theCmd)

  if !reflect.DeepEqual(*command, theCmd) {
    t.Error("DeepEqual failed.")
  }
}

func TestTextToCommand_0Commands_2_KV_White(t *testing.T) {
  t.Log("TextToCommand: 0 commands, 2 KV with white space")
  sReq := Request {
    Token: "abcd1234",
    TeamDomain: "example",
    ChannelId: "C2147483705",
    ChannelName: "test",
    UserId: "U2147483697",
    UserName: "Steve",
    Command: "/devhub",
    Text: " \n  \n\n title \n  \n\n = \n  \n\n test number 7 \n  \n\n labels \n  \n\n = \n  \n\n EPS, somethingElse \n  \n\n ",
    ResponseUrl: "http://localhost:8002/test",
  }

  command, err := sReq.TextToCommand()
  t.Logf("Request.TextToCommand: %#v", command)
  if err != nil {
    t.Errorf("Epic fail!! %v", err)
    return
  }

  theCmd := DevHubCommand {
    []string(nil),
    map[string]string{
      "title": "test number 7",
      "labels": "EPS, somethingElse",
    },
  }
  t.Logf("DevHubCommand: %#v", theCmd)

  if !reflect.DeepEqual(*command, theCmd) {
    t.Error("DeepEqual failed.")
  }
}
