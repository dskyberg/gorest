// Package config is just a simple wrapper around github.com/spf13/viper
package config

import (
  "fmt"
  "log"
  "bytes"
  "time"
  "strings"
  "path/filepath"
  "github.com/fsnotify/fsnotify"
  "github.com/spf13/cast"
  "github.com/spf13/viper"
)

// Config is a wrapper for Viper. You can call any Viper method directly on
// the Config object.  So in:
//     c := config.New(...)
//     c.Get("value")
// Viper.Get is actually called.
type Config struct {
  viper.Viper
}

// Creates a new Config/Viper instance, and reads the config file.
// The fileName should NOT have an extension.  Viper automatically
// searches for the files of the types it supports.  Which includes
// json, yaml, toml, props, properties, etc..
func New(fileName *string, defaults *map[string]interface{}) *Config {

  c := Config{*viper.New()}

  // Always look in the ENV for values, as well as the config file
  c.AutomaticEnv()

  // Load the defaults
  if defaults != nil {
    for key, value := range *defaults {
      c.SetDefault(key, value)
    }
  }

  // Set the config file path and name (if provided)
  var (
    path string
    config string
  )

  if fileName != nil {
    path, config = ParseConfig(*fileName)
  } else {
    path, config = ParseConfig("")
  }
  c.AddConfigPath(".")
  // Always add the current folder, unless already specified
  if path != "." {
    c.AddConfigPath(".")
  }
  if (config != "") {
    c.SetConfigName(config)
  }

  log.Printf("Loading config file: [%s]  from [%s]", config, path)
  err := c.ReadInConfig() // Find and read the config file
  if err != nil { // Handle errors reading the config file
    log.Printf("Error reading config file: %s \n", err)
  } else {
    c.WatchConfig()
    c.OnConfigChange(func(e fsnotify.Event) {
      log.Println("Config file changed:", e.Name)
    })
  }

  return &c
}

// Key is just a helper function that concats two strings much faster than
// Go's native +=
func Key(pre string, key string) string {
  var buffer bytes.Buffer
  buffer.WriteString(pre)
  buffer.WriteString(key)
  return buffer.String()
}

func Join(sep string, args ...string) string {
  var buffer bytes.Buffer
  numArgs := len(args)
  for i, part := range args {
    buffer.WriteString(part)
    if len(sep) > 0 && i < numArgs - 1 {
      buffer.WriteString(sep)
    }
  }
  return buffer.String()
}
// ParseConfig breaks the provided fileName into a path and base file name.
// If the fileName is "", ".", or "..", the the path will be "." and the
// config will be ""
// See Config_internal_ParseConfig_test.go for examples
func ParseConfig(fileNameInput string) (string, string) {
  fmt.Printf("\nParseConfig: fileNameInput [%s]\n", fileNameInput)
  if len(fileNameInput) == 0 {
    return ".", ""
  }

  fileName := filepath.Clean(fileNameInput)
  fmt.Printf("ParseConfig: fileName [%s]\n", fileName)
  if fileName == "." || fileName == ".." {
    return fileName, ""
  }

  if fileNameInput[len(fileNameInput)-1] == filepath.Separator {
    return fileName, ""
  }

  // If there is no path separator, then this is just
  // Dir will return "." for empty fileName
  path := filepath.Dir(fileName)
  fmt.Printf("ParseConfig: path [%s]\n",path)


  var config string
  // Short circuit for empty ( == ".") path
  base := filepath.Base(fileName)
  fmt.Printf("ParseConfig: base [%s]\n",base)
  if base == "." || base == ".."{
    config = ""
  } else {
    extIdx := strings.Index(base, ".")
    if extIdx > 0 {
      parts := strings.Split(base, ".")
      config = parts[0]
    } else {
      config = base
    }
  }
  fmt.Printf("ParseConfig: config [%s]\n", config)
  return path, config
}

// GetOrDefault looks for the key in Viper and returns, if found.
// Else, def is returned.
func (c *Config) GetOrDefault(key string, def interface{} ) interface{} {
  if i := c.Get(key); i != nil {
    return i
  } else {
    return def
  }
}

// GetBoolOrDefault looks for the key in Viper and returns, if found.
// Else, def is returned.
func (c *Config) GetBoolOrDefault(key string, def bool ) bool {
  if i := c.Get(key); i != nil {
    return cast.ToBool(i)
  } else {
    return def
  }
}

// GetFloat64OrDefault looks for the key in Viper and returns, if found.
// Else, def is returned.
func (c *Config) GetFloat64OrDefault(key string, def float64 ) float64 {
  if i := c.Get(key); i != nil {
    return cast.ToFloat64(i)
  } else {
    return def
  }
}

// GetIntOrDefault looks for the key in Viper and returns, if found.
// Else, def is returned.
func (c *Config) GetIntOrDefault(key string, def int ) int {
  if i := c.Get(key); i != nil {
    return cast.ToInt(i)
  } else {
    return def
  }
}

// GetStringOrDefault looks for the key in Viper and returns, if found.
// Else, def is returned.
func (c *Config) GetStringOrDefault(key string, def string ) string {
  if i := c.Get(key); i != nil {
    return cast.ToString(i)
  } else {
    return def
  }
}

// GetStringMapOrDefault looks for the key in Viper and returns, if found.
// Else, def is returned.
func (c *Config) GetStringMapOrDefault(key string, def map[string]interface{} ) map[string]interface{} {
  if i := c.Get(key); i != nil {
    return cast.ToStringMap(i)
  } else {
    return def
  }
}

// GetStringMapStringOrDefault looks for the key in Viper and returns, if found.
// Else, def is returned.
func (c *Config) GetStringMapStringOrDefault(key string, def map[string]string ) map[string]string {
  if i := c.Get(key); i != nil {
    return cast.ToStringMapString(i)
  } else {
    return def
  }
}

// GetStringSliceOrDefault looks for the key in Viper and returns, if found.
// Else, def is returned.
func (c *Config) GetStringSliceOrDefault(key string, def []string ) []string {
  if i := c.Get(key); i != nil {
    return cast.ToStringSlice(i)
  } else {
    return def
  }
}

// GetTimeOrDefault looks for the key in Viper and returns, if found.
// Else, def is returned.
func (c *Config) GetTimeOrDefault(key string, def time.Time ) time.Time {
  if i := c.Get(key); i != nil {
    return cast.ToTime(i)
  } else {
    return def
  }
}

// GetDurationOrDefault looks for the key in Viper and returns, if found.
// Else, def is returned.
func (c *Config) GetDurationOrDefault(key string, def time.Duration ) time.Duration {
  if i := c.Get(key); i != nil {
    return cast.ToDuration(i)
  } else {
    return def
  }
}
