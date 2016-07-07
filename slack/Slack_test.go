package slack

import (
  "testing"
  //"reflect"
  "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/require"
)

type RequestCommand struct {
  Request Request
  Command DevHubCommand
  KvText string
}


var requestCommands = map[string]RequestCommand {
  "0 commands, 0 KV": RequestCommand {
    Request {
      Token: "abcd1234",
      TeamDomain: "example",
      ChannelId: "C2147483705",
      ChannelName: "test",
      UserId: "U2147483697",
      UserName: "Steve",
      Command: "/devhub",
      Text: "",
      ResponseUrl: "http://localhost:8002/test",
    },
    DevHubCommand {
      []string{},
      map[string]string{},
    },
    "",
  },
  "0 commands, 0 KV with white space": RequestCommand {
    Request {
      Token: "abcd1234",
      TeamDomain: "example",
      ChannelId: "C2147483705",
      ChannelName: "test",
      UserId: "U2147483697",
      UserName: "Steve",
      Command: "/devhub",
      Text: "  \n \n\n ",
      ResponseUrl: "http://localhost:8002/test",
    },
    DevHubCommand {
      []string{},
      map[string]string{},
    },
    "",
  },
  "1 commands, 0 KV": RequestCommand {
    Request {
      Token: "abcd1234",
      TeamDomain: "example",
      ChannelId: "C2147483705",
      ChannelName: "test",
      UserId: "U2147483697",
      UserName: "Steve",
      Command: "/devhub",
      Text: "cmd1",
      ResponseUrl: "http://localhost:8002/test",
    },
    DevHubCommand {
      []string{"cmd1"},
      map[string]string{},
    },
    "",
  },
  "1 commands, 0 KV with white space": RequestCommand {
    Request {
      Token: "abcd1234",
      TeamDomain: "example",
      ChannelId: "C2147483705",
      ChannelName: "test",
      UserId: "U2147483697",
      UserName: "Steve",
      Command: "/devhub",
      Text: "  \n \r\n cmd1  \n \r\n ",
      ResponseUrl: "http://localhost:8002/test",
    },
    DevHubCommand {
      []string{"cmd1"},
      map[string]string{},
    },
    "",
  },
  "2 commands, 0 KV": RequestCommand {
    Request {
      Token: "abcd1234",
      TeamDomain: "example",
      ChannelId: "C2147483705",
      ChannelName: "test",
      UserId: "U2147483697",
      UserName: "Steve",
      Command: "/devhub",
      Text: "cmd1 cmd2",
      ResponseUrl: "http://localhost:8002/test",
    },
    DevHubCommand {
      []string{"cmd1", "cmd2"},
      map[string]string{},
    },
    "",
  },
  "2 commands, 0 KV with white space": RequestCommand {
    Request {
      Token: "abcd1234",
      TeamDomain: "example",
      ChannelId: "C2147483705",
      ChannelName: "test",
      UserId: "U2147483697",
      UserName: "Steve",
      Command: "/devhub",
      Text: "  \n \n\n cmd1  \n \r\n   \n \r\n cmd2  \n \r\n ",
      ResponseUrl: "http://localhost:8002/test",
    },
    DevHubCommand {
      []string{"cmd1", "cmd2"},
      map[string]string{},
    },
    "",
  },
  "0 commands, 1 KV": RequestCommand {
    Request {
      Token: "abcd1234",
      TeamDomain: "example",
      ChannelId: "C2147483705",
      ChannelName: "test",
      UserId: "U2147483697",
      UserName: "Steve",
      Command: "/devhub",
      Text: "title=test number 7",
      ResponseUrl: "http://localhost:8002/test",
    },
    DevHubCommand {
      []string{},
      map[string]string {
        "title": "test number 7",
      },
    },
    "title=test number 7",
  },
  "0 commands, 1 KV with white space": RequestCommand {
    Request {
      Token: "abcd1234",
      TeamDomain: "example",
      ChannelId: "C2147483705",
      ChannelName: "test",
      UserId: "U2147483697",
      UserName: "Steve",
      Command: "/devhub",
      Text: " \n  \n\n title \n  \n\n = \n  \n\n test number 7 \n  \n\n ",
      ResponseUrl: "http://localhost:8002/test",
    },
    DevHubCommand {
      []string{},
      map[string]string {
        "title": "test number 7",
      },
    },
    " \n  \n\n title \n  \n\n = \n  \n\n test number 7 \n  \n\n ",
  },
  "0 commands, 2 KV": RequestCommand {
    Request {
      Token: "abcd1234",
      TeamDomain: "example",
      ChannelId: "C2147483705",
      ChannelName: "test",
      UserId: "U2147483697",
      UserName: "Steve",
      Command: "/devhub",
      Text: "title=test number 7 labels=EPS, somethingElse",
      ResponseUrl: "http://localhost:8002/test",
    },
    DevHubCommand {
      []string{},
      map[string]string{
        "title": "test number 7",
        "labels": "EPS, somethingElse",
      },
    },
    "title=test number 7 labels=EPS, somethingElse",
  },
  "0 commands, 2 KV with white space": RequestCommand {
    Request {
      Token: "abcd1234",
      TeamDomain: "example",
      ChannelId: "C2147483705",
      ChannelName: "test",
      UserId: "U2147483697",
      UserName: "Steve",
      Command: "/devhub",
      Text: " \n  \n\n title \n  \n\n = \n  \n\n test number 7 \n  \n\n labels \n  \n\n = \n  \n\n EPS, somethingElse \n  \n\n ",
      ResponseUrl: "http://localhost:8002/test",
    },
    DevHubCommand {
      []string{},
      map[string]string{
        "title": "test number 7",
        "labels": "EPS, somethingElse",
      },
    },
    " \n  \n\n title \n  \n\n = \n  \n\n test number 7 \n  \n\n labels \n  \n\n = \n  \n\n EPS, somethingElse \n  \n\n ",
  },
}

func TestParseCommands(t *testing.T) {
  for name, rq := range requestCommands {
    t.Logf("ParseCommands: %s", name)

    commands, kvText := ParseCommands(rq.Request.Text)
    assert.Equal(t, rq.Command.Commands, commands,
      "ParseCommands returned the wrong number of commands: %#v", commands)
    assert.Equal(t, rq.KvText, kvText,
      "ParseCommands returned the wrong remainder: %s", kvText)
  }
}

func TestParseKeyValuePairs(t *testing.T) {
  for name, rq := range requestCommands {
    t.Log("ParseKeyValuePairs: %s", name)
    kv, err := ParseKeyValuePairs(rq.Request.Text)
    require.Nil(t, err, "ParseKeyValuePairs returned an error: %#v", err)
    assert.Equal(t, rq.Command.Params, kv,
      "ParseKeyValuePairs returned a bad set of kvs: %#v", kv)
  }
}

func TestTextToCommand(t *testing.T) {
  for name, rq := range requestCommands {
    t.Logf("TextToCommand: %s", name)
    command, err := rq.Request.TextToCommand()
    require.Nil(t, err, "Epic fail!! %v", err)
    assert.Equal(t, rq.Command, *command, "%s failed.", name)
  }
}
