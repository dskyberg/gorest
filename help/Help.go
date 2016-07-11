// Package help is a utility package that enables expressing
// help text in an HCL file.  HCL is used because of it's ability
// handle richer, multiline values.
package help

import (
  "os"
  "io"
  "io/ioutil"
  "strings"
  "fmt"
  "errors"
  "github.com/hashicorp/hcl"

)

type Help map[string]string

const Base = "base"
const Sep = "."

func (h Help) Base() string {
  return h.Get(Base)
}

func (h Help) Get(key string) string {
  if text, ok := h[key]; ok {
    return text
  }
  return h[Base]
}

func ParseHelpText(helpText string) (Help, error) {
  reader := strings.NewReader(helpText)
  if reader == nil {
    return nil, errors.New("Could not load text from provided string")
  }
  return ParseHelp(reader)
}


func ParseHelpFile(fileName string) (Help, error) {

  reader, err := os.Open(fileName)
  if err != nil {
    return nil, errors.New(fmt.Sprintf("Could open help file %s: %s",
      fileName, err.Error()))
  }
  return ParseHelp(reader)
}

// ParseHelpText
func ParseHelp(reader io.Reader) (Help, error) {

  helpText, err := ioutil.ReadAll(reader)
  if err != nil {
    return nil, errors.New("Could not load help")
  }
  var help Help
  err = hcl.Decode(&help, string(helpText))
  if err != nil {
    return nil, errors.New(fmt.Sprintf("Could not parse help text", err.Error()))
  }
  if help == nil || len(help) == 0 {
    return nil, errors.New("Empty help file")
  }

  // Make sure there is a base help
  if _, ok := help[Base]; !ok {
    return nil, errors.New("Help does not containe a base help")
  }

  return help, nil
}
