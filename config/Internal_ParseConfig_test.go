package config

import (
  "fmt"
  "testing"

  . "github.com/smartystreets/goconvey/convey"
)

func TestParseConfig(t *testing.T) {
  var (
    filePath string
    expectedPath string
    expectedConfig string
    path string
    config string
  )

  filePath = ""
  expectedPath = "."
  expectedConfig = ""
  path, config = ParseConfig(filePath)
  Convey(fmt.Sprintf("Given the path [%s]", filePath), t, func() {
    Convey(fmt.Sprintf("The path should be [%s]", expectedPath), func() {
      So(path, ShouldEqual, expectedPath)
    })
    Convey(fmt.Sprintf("And the config should be [%s]", expectedConfig), func() {
      So(config, ShouldEqual, expectedConfig)
    })
  })

  filePath = "."
  expectedPath = "."
  expectedConfig = ""
  path, config = ParseConfig(filePath)
  Convey(fmt.Sprintf("Given the path [%s]", filePath), t, func() {
    Convey(fmt.Sprintf("The path should be [%s]", expectedPath), func() {
      So(path, ShouldEqual, expectedPath)
    })
    Convey(fmt.Sprintf("And the config should be [%s]", expectedConfig), func() {
      So(config, ShouldEqual, expectedConfig)
    })
  })

  filePath = "./"
  expectedPath = "."
  expectedConfig = ""
  path, config = ParseConfig(filePath)
  Convey(fmt.Sprintf("Given the path [%s]", filePath), t, func() {
    Convey(fmt.Sprintf("The path should be [%s]", expectedPath), func() {
      So(path, ShouldEqual, expectedPath)
    })
    Convey(fmt.Sprintf("And the config should be [%s]", expectedConfig), func() {
      So(config, ShouldEqual, expectedConfig)
    })
  })

  filePath = ".."
  expectedPath = ".."
  expectedConfig = ""
  path, config = ParseConfig(filePath)
  Convey(fmt.Sprintf("Given the path [%s]", filePath), t, func() {
    Convey(fmt.Sprintf("The path should be [%s]", expectedPath), func() {
      So(path, ShouldEqual, expectedPath)
    })
    Convey(fmt.Sprintf("And the config should be [%s]", expectedConfig), func() {
      So(config, ShouldEqual, expectedConfig)
    })
  })

  filePath = "../"
  expectedPath = ".."
  expectedConfig = ""
  path, config = ParseConfig(filePath)
  Convey(fmt.Sprintf("Given the path [%s]", filePath), t, func() {
    Convey(fmt.Sprintf("The path should be [%s]", expectedPath), func() {
      So(path, ShouldEqual, expectedPath)
    })
    Convey(fmt.Sprintf("And the config should be [%s]", expectedConfig), func() {
      So(config, ShouldEqual, expectedConfig)
    })
  })

  filePath = "part1"
  expectedPath = "."
  expectedConfig = "part1"
  path, config = ParseConfig(filePath)
  Convey(fmt.Sprintf("Given the path [%s]", filePath), t, func() {
    Convey(fmt.Sprintf("The path should be [%s]", expectedPath), func() {
      So(path, ShouldEqual, expectedPath)
    })
    Convey(fmt.Sprintf("And the config should be [%s]", expectedConfig), func() {
      So(config, ShouldEqual, expectedConfig)
    })
  })

  filePath = "part1/"
  expectedPath = "part1"
  expectedConfig = ""
  path, config = ParseConfig(filePath)
  Convey(fmt.Sprintf("Given the path [%s]", filePath), t, func() {
    Convey(fmt.Sprintf("The path should be [%s]", expectedPath), func() {
      So(path, ShouldEqual, expectedPath)
    })
    Convey(fmt.Sprintf("And the config should be [%s]", expectedConfig), func() {
      So(config, ShouldEqual, expectedConfig)
    })
  })

  filePath = "part1/."
  expectedPath = "."
  expectedConfig = "part1"
  path, config = ParseConfig(filePath)
  Convey(fmt.Sprintf("Given the path [%s]", filePath), t, func() {
    Convey(fmt.Sprintf("The path should be [%s]", expectedPath), func() {
      So(path, ShouldEqual, expectedPath)
    })
    Convey(fmt.Sprintf("And the config should be [%s]", expectedConfig), func() {
      So(config, ShouldEqual, expectedConfig)
    })
  })

  filePath = "/part1"
  expectedPath = "/"
  expectedConfig = "part1"
  path, config = ParseConfig(filePath)
  Convey(fmt.Sprintf("Given the path [%s]", filePath), t, func() {
    Convey(fmt.Sprintf("The path should be [%s]", expectedPath), func() {
      So(path, ShouldEqual, expectedPath)
    })
    Convey(fmt.Sprintf("And the config should be [%s]", expectedConfig), func() {
      So(config, ShouldEqual, expectedConfig)
    })
  })

  filePath = "part1/part2"
  expectedPath = "part1"
  expectedConfig = "part2"
  path, config = ParseConfig(filePath)
  Convey(fmt.Sprintf("Given the path [%s]", filePath), t, func() {
    Convey(fmt.Sprintf("The path should be [%s]", expectedPath), func() {
      So(path, ShouldEqual, expectedPath)
    })
    Convey(fmt.Sprintf("And the config should be [%s]", expectedConfig), func() {
      So(config, ShouldEqual, expectedConfig)
    })
  })

  filePath = "part1/part2/"
  expectedPath = "part1/part2"
  expectedConfig = ""
  path, config = ParseConfig(filePath)
  Convey(fmt.Sprintf("Given the path [%s]", filePath), t, func() {
    Convey(fmt.Sprintf("The path should be [%s]", expectedPath), func() {
      So(path, ShouldEqual, expectedPath)
    })
    Convey(fmt.Sprintf("And the config should be [%s]", expectedConfig), func() {
      So(config, ShouldEqual, expectedConfig)
    })
  })

  filePath = "part1.ext"
  expectedPath = "."
  expectedConfig = "part1"
  path, config = ParseConfig(filePath)
  Convey(fmt.Sprintf("Given the path [%s]", filePath), t, func() {
    Convey(fmt.Sprintf("The path should be [%s]", expectedPath), func() {
      So(path, ShouldEqual, expectedPath)
    })
    Convey(fmt.Sprintf("And the config should be [%s]", expectedConfig), func() {
      So(config, ShouldEqual, expectedConfig)
    })
  })

  filePath = "./part1.ext"
  expectedPath = "."
  expectedConfig = "part1"
  path, config = ParseConfig(filePath)
  Convey(fmt.Sprintf("Given the path [%s]", filePath), t, func() {
    Convey(fmt.Sprintf("The path should be [%s]", expectedPath), func() {
      So(path, ShouldEqual, expectedPath)
    })
    Convey(fmt.Sprintf("And the config should be [%s]", expectedConfig), func() {
      So(config, ShouldEqual, expectedConfig)
    })
  })

  filePath = "../part1.ext"
  expectedPath = ".."
  expectedConfig = "part1"
  path, config = ParseConfig(filePath)
  Convey(fmt.Sprintf("Given the path [%s]", filePath), t, func() {
    Convey(fmt.Sprintf("The path should be [%s]", expectedPath), func() {
      So(path, ShouldEqual, expectedPath)
    })
    Convey(fmt.Sprintf("And the config should be [%s]", expectedConfig), func() {
      So(config, ShouldEqual, expectedConfig)
    })
  })

  filePath = "part1/part2.ext"
  expectedPath = "part1"
  expectedConfig = "part2"
  path, config = ParseConfig(filePath)
  Convey(fmt.Sprintf("Given the path [%s]", filePath), t, func() {
    Convey(fmt.Sprintf("The path should be [%s]", expectedPath), func() {
      So(path, ShouldEqual, expectedPath)
    })
    Convey(fmt.Sprintf("And the config should be [%s]", expectedConfig), func() {
      So(config, ShouldEqual, expectedConfig)
    })
  })

  filePath = "./part1/part2.ext"
  expectedPath = "part1"
  expectedConfig = "part2"
  path, config = ParseConfig(filePath)
  Convey(fmt.Sprintf("Given the path [%s]", filePath), t, func() {
    Convey(fmt.Sprintf("The path should be [%s]", expectedPath), func() {
      So(path, ShouldEqual, expectedPath)
    })
    Convey(fmt.Sprintf("And the config should be [%s]", expectedConfig), func() {
      So(config, ShouldEqual, expectedConfig)
    })
  })

  filePath = "/part1/part2.ext"
  expectedPath = "/part1"
  expectedConfig = "part2"
  path, config = ParseConfig(filePath)
  Convey(fmt.Sprintf("Given the path [%s]", filePath), t, func() {
    Convey(fmt.Sprintf("The path should be [%s]", expectedPath), func() {
      So(path, ShouldEqual, expectedPath)
    })
    Convey(fmt.Sprintf("And the config should be [%s]", expectedConfig), func() {
      So(config, ShouldEqual, expectedConfig)
    })
  })

}
