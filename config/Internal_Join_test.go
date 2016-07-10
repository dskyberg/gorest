package config

import (
  "fmt"
  "testing"

  . "github.com/smartystreets/goconvey/convey"
)

var NUMBER_SETTING_1 = 0
var NUMBER_SETTING_2 = 1
var STRING_SETTING_1 = "Test String 1"
var STRING_SETTING_2 = "Test String 2"
var configDefaults = map[string]interface{} {
  "NUMBER_SETTING_1": NUMBER_SETTING_1,
  "NUMBER_SETTING_2": NUMBER_SETTING_2,
  "STRING_SETTING_1": STRING_SETTING_1,
  "STRING_SETTING_2": STRING_SETTING_2,
}
var defaultConfigName = "config"
var configFile = ".."
/*
// ConfigTestSuite extends (embeds) Suite, and allows us to
// create setup data for the test run.
type ConfigTestSuite struct {
  suite.Suite
  config *Config
}

// Official Test<Method> for go test to find.  This just loads and
// runs our defined suite of tests, which follow.
func TestConfigTestSuite(t* testing.T) {
  suite.Run(t, new(ConfigTestSuite))
}

// Suite setup routine.  Load the ConfigTestSuite as needed for the test
func (suite *ConfigTestSuite) SetupTest() {

  suite.config = New(configFile, configDefaults)
}

func (suite *ConfigTestSuite) TesJoin() {
  data := []struct {
    Sep string
    Expected string
  } {
    { "", "Val1Val2VAl3"
  }
  sep1 := ""
  sep2 := " "
  sep3 := "_"
  sep4 := "_tst_"

  joined1 := "Val1Val2VAl3"
  joined2 := "Val1_Val2_Val3"
  joined3 := "Val1_tst_Val2_tst_Val3"

  // Empty separator
  assert.Equal(suite.T(), joined1, Join(sep1, "Val1", "Val2", "Val3"), "Wrong joined response")
  assert.Equal(suite.T(), joined1, Join(sep1, "", "Val1", "Val2", "Val3", ""), "Wrong joined response")

  assert.Equal(suite.T(), joined2, Join(sep2, "Val1", "Val2", "Val3"), "Wrong joined response")
  assert.Equal(suite.T(), joined2, Join(sep2, "Val1", "Val2", "Val3", ""), "Wrong joined response")

  assert.Equal(suite.T(), joined4, Join(sep4, "Val1", "Val2", "Val3"), "Wrong joined response")
}


// Test to see if Config defaults are properly configured.
func (suite *ConfigTestSuite) TestDefaults() {
  assert.Equal(suite.T(), NUMBER_SETTING_1, suite.config.GetInt("NUMBER_SETTING_1"), "Wrong default")
}
*/


func TestJoin(t *testing.T) {
  Convey("Given the separator \"\"", t, func() {
    sep := ""

    Convey("When the values are Val1, Val2, Val3", func() {
      expectedResult := "Val1Val2Val3"
      Convey(fmt.Sprintf("The joined value should be %s", expectedResult), func() {
        result := Join(sep, "Val1", "Val2", "Val3")
        So(result, ShouldEqual, expectedResult)
      })
    })

    Convey("When the values are \"\", Val1, Val2, Val3, \"\"", func() {
      expectedResult := "Val1Val2Val3"
      Convey(fmt.Sprintf("The joined value should be %s", expectedResult), func() {
        result := Join(sep, "", "Val1", "Val2", "Val3", "")
        So(result, ShouldEqual, result)
      })
    })

  })

  Convey("Given the separator \"_\"", t, func() {
    sep := "_"

    Convey("When the values are Val1, Val2, Val3", func() {
      expectedResult := "Val1_Val2_Val3"
      Convey(fmt.Sprintf("The joined value should be %s", expectedResult), func() {
        result := Join(sep, "Val1", "Val2", "Val3")
        So(result, ShouldEqual, expectedResult)
      })
    })

    Convey("When the values are \"\", Val1, Val2, Val3, \"\"", func() {
      expectedResult := "_Val1_Val2_Val3_"
      Convey(fmt.Sprintf("The joined value should be %s", expectedResult), func() {
        result := Join(sep, "", "Val1", "Val2", "Val3", "")
        So(result, ShouldEqual, expectedResult)
      })
    })

  })

  Convey("Given the separator \"_\"", t, func() {
    sep := "_tst_"

    Convey("When the values are Val1, Val2, Val3", func() {
      expectedResult := "Val1_tst_Val2_tst_Val3"
      Convey(fmt.Sprintf("The joined value should be %s", expectedResult), func() {
        result := Join(sep, "Val1", "Val2", "Val3")
        So(result, ShouldEqual, expectedResult)
      })
    })

    Convey("When the values are \"\", Val1, Val2, Val3, \"\"", func() {
      expectedResult := "_tst_Val1_tst_Val2_tst_Val3_tst_"
      result := Join(sep, "", "Val1", "Val2", "Val3", "")
      Convey(fmt.Sprintf("The joined value should be %s", expectedResult), func() {
        So(result, ShouldEqual, expectedResult)
      })
    })

  })

}
