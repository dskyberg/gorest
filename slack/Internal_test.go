package slack

import (
  "fmt"
  "testing"

  . "github.com/smartystreets/goconvey/convey"
)

type RequestCommand struct {
  Request Request
  Command DevHubCommand
  KvText string
}


var requestCommands = []RequestCommand {
/*
  RequestCommand {
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
      Commands{},
      map[string]string{},
    },
    "",
  },

  RequestCommand {
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
      Commands{},
      map[string]string{},
    },
    "",
  },
  RequestCommand {
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
      Commands{"cmd1"},
      map[string]string{},
    },
    "",
  },
  RequestCommand {
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
      Commands{"cmd1"},
      map[string]string{},
    },
    "",
  },
  RequestCommand {
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
      Commands{"cmd1", "cmd2"},
      map[string]string{},
    },
    "",
  },
  RequestCommand {
    Request {
      Token: "abcd1234",
      TeamDomain: "example",
      ChannelId: "C2147483705",
      ChannelName: "test",
      UserId: "U2147483697",
      UserName: "Steve",
      Command: "/devhub",
      Text: "cmd1\r\ncmd2",
      ResponseUrl: "http://localhost:8002/test",
    },
    DevHubCommand {
      Commands{"cmd1", "cmd2"},
      map[string]string{},
    },
    "",
  },
  RequestCommand {
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
      Commands{},
      map[string]string {
        "title": "test number 7",
      },
    },
    "title=test number 7",
  },
  RequestCommand {
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
      Commands{},
      map[string]string {
        "title": "test number 7",
      },
    },
    " \n  \n\n title \n  \n\n = \n  \n\n test number 7 \n  \n\n ",
  },
  RequestCommand {
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
      Commands{},
      map[string]string{
        "title": "test number 7",
        "labels": "EPS, somethingElse",
      },
    },
    "title=test number 7 labels=EPS, somethingElse",
  },
  */
  RequestCommand {
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
      Commands{},
      map[string]string{
        "title": "test number 7",
        "labels": "EPS, somethingElse",
      },
    },
    " \n  \n\n title \n  \n\n = \n  \n\n test number 7 \n  \n\n labels \n  \n\n = \n  \n\n EPS, somethingElse \n  \n\n ",
  },
}

func TestParseCommands(t *testing.T) {
  for _, rq := range requestCommands {
    commands, kvText := ParseCommands(rq.Request.Text)
    kv, err := ParseKeyValuePairs(rq.Request.Text)
    Convey(fmt.Sprintf("Given the input [%v]", rq.Request.Text), t, func() {
      Convey(fmt.Sprintf("The commands should be [%v]", rq.Command.Commands), func() {
        So(commands, ShouldResemble, rq.Command.Commands)
      })
      Convey(fmt.Sprintf("The KV Text should be [%v]", rq.KvText), func() {
        So(kvText, ShouldEqual, rq.KvText)
      })
      Convey("err should be nil", func() {
        So(err, ShouldBeNil)
      })
      Convey(fmt.Sprintf("The KV Pairs should be [%v]", rq.Command.Params), func() {
        So(kv, ShouldEqual, rq.Command.Params)
      })
    })
  }
}

func TestParseKeyValuePairs(t *testing.T) {
  for _, rq := range requestCommands {
    kv, err := ParseKeyValuePairs(rq.Request.Text)
    Convey(fmt.Sprintf("Given the input [%v]", rq.Request.Text), t, func() {
      Convey("Err should be nil", func() {
        So(err, ShouldBeNil)
      })
      Convey(fmt.Sprintf("The KV Pairs should be [%v]", rq.Command.Params), func() {
        So(kv, ShouldResemble, rq.Command.Params)
      })
    })
  }
}

/*
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
*/
