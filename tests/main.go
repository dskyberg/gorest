package main


import (
  "fmt"
  "reflect"
  "path/filepath"
  "strings"
  "github.com/confyrm/gorest/app/cmd"
)

func main() {
  // Run the admin port
  //BasicTest()
  FilePathTest()
}


func FilePathTest() {
  path := "/Users/david/golang/src/github.com/confyrm/gorest/config.yml"
  base := filepath.Base(path)
  parts := strings.Split(base, ".")
  fmt.Printf("Config path: %s\nConfig file: %s\n", filepath.Dir(path), parts[0])
}

func BasicTest() {
  cs := []string {"cmd1","cmd2"}
  kv := make(map[string]string)
  kv["title"] = "test number 7"
  kv["labels"] = "EPS, otherLabel"
  theCmd := cmd.DevHubCommand{cs, kv}
  _ = theCmd
  text := "cmd1    cmd2    title  =   Test number 7 labels   = EPS, otherLabel"

  command, err := cmd.TextToCommand(text)
  if err != nil {
    fmt.Printf("Epic fail!! %v", err)
    return
  }
  if !reflect.DeepEqual(*command, theCmd) {
    fmt.Printf("DeepEqual failed\n")
    fmt.Printf("Commands: %#v\n", command.Commands)
    fmt.Printf("Params: %#v\n", command.Params)
  }
  fmt.Printf("DeepEqual Succeeded:\n")
}

func NoKVTest() {
  cs := []string {"cmd1","cmd2"}
  theCmd := cmd.DevHubCommand{cs, nil}
  _ = theCmd
  text := "cmd1    cmd2  "

  command, err := cmd.TextToCommand(text)
  if err != nil {
    fmt.Printf("Epic fail!! %v", err)
    return
  }
  if !reflect.DeepEqual(command, theCmd) {
    fmt.Printf("DeepEqual failed\n")
    fmt.Printf("Commands: %#v\n", command.Commands)
    fmt.Printf("Params: %#v\n", command.Params)
  }
  fmt.Printf("DeepEqual Succeeded:\n")
}
