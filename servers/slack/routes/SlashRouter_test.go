package routes

import (
  "testing"

  "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/require"
)
const helpPath = "../../../help/help.hcl"
const testKey = "test.some.stuff"
const testValue = "Help for test some stuff"

func TestParseHelpText(t *testing.T) {
  out, err := ParseHelpText(helpPath)
  require.Nil(t, err, "Epic fail!! ParseHelpText returned an error: %v", err)

  t.Logf("out: %#v", out)

  txt := out["test.some.stuff"];

  assert.Equal(t, testValue, txt.(string),
    "Epic fail!! ParseHelpText: value found for [%s] %s", testKey, txt)
}
