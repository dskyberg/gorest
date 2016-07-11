package config

import (
  "fmt"
  "testing"

  . "github.com/smartystreets/goconvey/convey"
)

var testData = []struct {
  FilePath string
  ExpectedPath string
  ExpectedConfig string
} {
  {
    "",
    ".",
    "",
  },
  {
    ".",
    ".",
    "",
  },
  {
    "./",
    ".",
    "",
  },
  {
    "..",
    "..",
    "",
  },
  {
    "../",
    "..",
    "",
  },
  {
    "part1",
    ".",
    "part1",
  },
  {
    "part1/",
    "part1",
    "",
  },
  {
    "part1/.",
    ".",
    "part1",
  },
  {
    "/part1",
    "/",
    "part1",
  },
  {
    "part1/part2",
    "part1",
    "part2",
  },
  {
    "part1/part2/",
    "part1/part2",
    "",
  },
  {
    "part1.ext",
    ".",
    "part1",
  },
  {
    "./part1.ext",
    ".",
    "part1",
  },
  {
    "../part1.ext",
    "..",
    "part1",
  },
  {
    "part1/part2.ext",
    "part1",
    "part2",
  },
  {
    "./part1/part2.ext",
    "part1",
    "part2",
  },
  {
    "/part1/part2.ext",
    "/part1",
    "part2",
  },
}

func TestParseConfig(t *testing.T) {
  for _, d := range testData {
    path, config := ParseConfig(d.FilePath)
    Convey(fmt.Sprintf("Given the input [%s]", d.FilePath), t, func() {
      Convey(fmt.Sprintf("The path should be [%s],", d.ExpectedPath), func() {
        So(path, ShouldEqual, d.ExpectedPath)
      })
      Convey(fmt.Sprintf("And the config should be [%s].", d.ExpectedConfig), func() {
        So(config, ShouldEqual, d.ExpectedConfig)
      })
    })
  }
}
