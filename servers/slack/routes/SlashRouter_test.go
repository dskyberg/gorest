package routes

import (
  "testing"
)

var testData = map[string]interface {} {
  "test.some.stuff": "Help for test some stuff",
}

func TestParseHelpText(t *testing.T) {
  out, err := ParseHelpText("../../../help/help.hcl")
  if err != nil {
    t.Errorf("Epic fail!! %v", err)
    return
  }
  t.Logf("out: %#v", out)
  var (
    txt interface{}
    ok bool
  )
  if txt, ok = out["test.some.stuff"]; !ok {
    t.Error("Epic fail!! No value found for [test.some.stuff]")
  }
  if txt != testData["test.some.stuff"]{
    t.Errorf("Epic fail!! Wrong value found for [test.some.stuff] %s", txt)
  }
}
