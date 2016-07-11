package slack

import (
  "fmt"
  "strings"
  "errors"

  "github.com/spf13/cast"
)

// The Slack slash commands that we process contain a consistent structure that
// can be leveraged for any type of activity.
type Commands []string

// Value returns the nth command as a string
func (commands Commands) Value(i int) (string, bool) {
  if len(commands) > i {
    return commands[i], true
  }
    return "", false
}

// Pop pulls the next command off the command stack.  If no commands
// exist, Pop returns false.
func (cPtr *Commands) Pop() (string, bool) {
  if cPtr == nil || len(*cPtr) == 0 {
    return "", false
  }
  res:=(*cPtr)[0]
  *cPtr=(*cPtr)[1:len(*cPtr)-1]
  return res, true
}

// Push puts a command on the command stack.
func (cPtr *Commands) Push(cmd string) {
  if cPtr == nil  {
    return
  }
  *cPtr = append([]string{cmd}, *cPtr...)
}

// Peek returns the next value, but leaves the value on the command stack
func (cPtr *Commands) Peek() (string, bool) {
  if cPtr == nil || len(*cPtr) == 0 {
    return "", false
  }
  return (*cPtr)[0], true
}

func (cPtr *Commands) HasHelp() bool {
  if cPtr == nil || len(*cPtr) == 0 {
    return false
  }
  for _, val := range *cPtr {
    if val == "help" {
      return true
    }
  }
  return false
}

// ValueToInt returns the nth command as an int
func (commands Commands) ValueToInt(i int) (int, bool) {
  if len(commands) > i {
    return cast.ToInt(commands[i]), true
  }
    return -1, false
}



type KVPairs map[string]string

type DevHubCommand struct {
  Commands Commands
  Params KVPairs
}

// IsEmpty returns true of both the Commands and the KVPairs are empty
func (command *DevHubCommand) IsEmpty() bool {
  if command == nil {
    return true
  }
  if len(command.Commands) == 0 && len(command.Params) == 0 {
    return true
  }
  return false
}

func (command *DevHubCommand) HasHelp() bool {
    if len(command.Params) == 0 {
      return command.Commands.HasHelp()
    }
    return false
}
func (command *DevHubCommand) HelpPath() (Commands, bool) {

  if command == nil {
    return Commands{}, true
  }

  cmdLen := len(command.Commands)
  if cmdLen == 0 {
    return Commands{}, true
  }

  lastIdx := cmdLen - 1
  head := command.Commands[0]
  tail := command.Commands[lastIdx]

  // No command, or single help command.
  // Like: /devhub or /devhub help
  if cmdLen == 0 || (cmdLen == 1 && head == "help") {
    return Commands{}, true
  }

  // Single command, with no params
  // like: /devhub new
  if cmdLen == 1 && len(command.Params) == 0 {
    return Commands{head}, true
  }

  // Command with "help" at the beginning or the end
  // Like /devhub help new
  // Like /devhub new help
  if head == "help" {
    return command.Commands[1:], true
  }
  if tail == "help" {
    return command.Commands[:lastIdx], true
  }

  // No help path
  return Commands{}, false
}

func (command *DevHubCommand) Pop() (string, bool) {
  return command.Commands.Pop()
}

func (command *DevHubCommand) Push(cmd string) {
  command.Commands.Push(cmd)
}

func (command *DevHubCommand) Peek() (string, bool) {
  return command.Commands.Peek()
}

func (command *DevHubCommand) CommandsFrom(from int) Commands {
    if from > len(command.Commands) {
      return Commands{}
    }
    return command.Commands[from:]
}

// ValueOrDefault returns either the value from the KV pairs,
// or the default, if the key is not found
func (command *DevHubCommand) ValueOrDefault(key string, def string) string {
  if value, ok := command.Params[key]; ok {
    return value
  } else {
    return def
  }
}

// Value is a simple wrapper around command.Params[key]
func (command *DevHubCommand) Value(key string) (string, bool) {
  value, ok := command.Params[key]
  return value, ok
}

// Value is a simple wrapper around command.Params[key]
func (command *DevHubCommand) Values(key string) ([]string, bool) {
  i, ok := command.Params[key]

  values := strings.Split(i, ",")
  for j := 0; j < len(values); j++ {
    values[j] = strings.TrimSpace(values[j])
  }
  return values, ok
}

// Value is a simple wrapper around command.Params[key]
func (command *DevHubCommand) ValueToInt(key string) (int, bool) {
  value, ok := command.Params[key]
  return cast.ToInt(value), ok
}

// Value is a simple wrapper around command.Params[key]
func (command *DevHubCommand) ValueToIntOrDefault(key string, def int) int {
  if value, ok := command.Params[key]; ok {
    return cast.ToInt(value)
  }
  return def
}

// HasValue is a utility function for command.Params
func (command *DevHubCommand) HasValue(key string) bool {
  if _, ok := command.Params[key]; ok {
    return true
  }
  return false
}



const TRIM_CUTSET = " \r\n"
const KV_DELIM = "="


func (sReq *Request) TextToCommand() (*DevHubCommand, error) {
  // Trim the string first, to remove any unwanted spaces and new lines
  t := strings.Trim(sReq.Text, TRIM_CUTSET)

  // Get the set of commands, and whatever text may be remaining after the commands
  commands, kvText := ParseCommands(t)
  kv, err := ParseKeyValuePairs(kvText)
  if err != nil {
    return nil, err
  }

  return &DevHubCommand{commands, kv}, nil
}

// ParseCommands is a helper function that parses out any commands that are
// placed before the Key/Value pairs in the provided text.
// Note, if there are no commands AND no KV pairs, ParseCommands will return
// an empty string for kvText.
func ParseCommands(text string) (Commands, string) {
  // Grab everything up to the first key
  var cmdText string
  var kvText string
  commands := make([]string, 0, 10)
  firstKey := FindStartOfNextKey(text)

  if firstKey == 0 {
    // There are KV's but no commands
    return commands, text
  } else if firstKey > 0 {
    cmdText = text[:firstKey-1]
    kvText = text[firstKey:]
  } else {
    // There doesn't appear to be any kv pairs.
    cmdText = text
    kvText = ""
  }
  // Read the list of commands first
  tmp := strings.Fields(cmdText)
  for i := 0; i < len(tmp); i++ {
    tmp[i] = strings.Trim(tmp[i], TRIM_CUTSET)
  }
  for _, x := range tmp {
    if len(x) > 0 {
      commands = append(commands, x)
    }
  }
  return commands, kvText
}

func ParseKeyValuePairs(kvText string) (KVPairs, error) {
  //Split on KV_DELIM.  This will yield an even number of values in s, where
  // s[i] = key, s[i+1] = value
  kv := make(KVPairs)

  for x := 0; x < 10; x++ { // Just to prevent infinite loop!
    i := strings.Index(kvText, KV_DELIM)

    if i == -1 {
      break
    }

    k := strings.ToLower(strings.Trim(kvText[:i], TRIM_CUTSET))
    //fmt.Printf(" - k: %s\n", k)
    kvText = strings.Trim(kvText[i + 1:], TRIM_CUTSET)

    if len(kvText) == 0 {
      err := errors.New(fmt.Sprintf("Error parsing Key Value pairs: No remainder for key: %s ", k ))
      return kv, err
    }

    j := FindStartOfNextKey(kvText)
    if j == -1 {
      // Last value in k/v pairs
      v := strings.Trim(kvText, TRIM_CUTSET)
      kv[k] = v
      //fmt.Printf(" - v: [%s]\n", v)
      break
    }


    // j now sits at letter of the next key.
    v := strings.ToLower(strings.Trim(kvText[:j - 1], TRIM_CUTSET))
    //fmt.Printf(" - v: [%s]\n", v)
    kv[k] = v
    kvText = strings.Trim(kvText[j:], TRIM_CUTSET)
  }
  return kv, nil
}


func FindStartOfNextKey(text string) int {
  t := strings.TrimLeft(text, TRIM_CUTSET)

  if len(t) == 0 {
    // No text to test
    return -1
  }
  i := strings.Index(t, KV_DELIM)

  if i == -1 {
    // No keys in the text
    return i
  }

  if i == 0 {
    // No keys, because text starts with '='
    return -1
  }

  if i == len(t) - 1 {
    // No keys, because text ends with '='
    return -1
  }

  // skip over any white space
  for i--; i > 0; i-- {
    if t[i] != ' ' && t[i] != '\n' && t[i] != '\r'  && t[i] != '\f' {
      break
    }
  }

  // Now skip the key
  for {
    if i == 0 {
      break
    }
    if t[i - 1] == ' ' {
      break
    }
    i--
  }
  return i
}
