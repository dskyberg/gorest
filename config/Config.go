// Package config is just a simple wrapper around github.com/spf13/viper
package config

import (
  "log"
  "bytes"
  "time"
  "strings"
  "path/filepath"
  "github.com/fsnotify/fsnotify"
  "github.com/spf13/cast"
  "github.com/spf13/viper"
)

// Config is a wrapper for Viper.
type Config struct {
  v *viper.Viper
}

// Creates a new Config/Viper instance, and reads the config file.
// The fileName should NOT have an extension.  Viper automatically
// searches for the files of the types it supports.  Which includes
// json, yaml, toml, props, properties, etc..
func New(fileName string, defaults map[string]interface{}) *Config {

  var path string
  var config string
  if fileName == "" {
    path = "."
    config = "config"
  } else {
    path, config = ParseConfig(fileName)
  }

  log.Printf("Loading config file: [%s]  from [%s]", config, path)

  v := viper.New()

  // Load the defaults
  for key, value := range defaults {
    v.SetDefault(key, value)
  }

  v.SetConfigName(config)

  v.AutomaticEnv()

  v.AddConfigPath(path)
  if path != "." {
    v.AddConfigPath(".")
  }
  err := v.ReadInConfig() // Find and read the config file
  if err != nil { // Handle errors reading the config file
    log.Printf("Error reading config file: %s \n", err)
  } else {
    v.WatchConfig()
    v.OnConfigChange(func(e fsnotify.Event) {
      log.Println("Config file changed:", e.Name)
    })
  }
  c := Config{v}
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

// ParseConfig breaks the provided fileName into a path, base name,
// and file type (based on the file ext, if present)
func ParseConfig(fileName string) (string, string) {

  path := filepath.Dir(fileName)
  base := filepath.Base(fileName)
  var config string
  extIdx := strings.Index(base, ".")
  if extIdx > 0 {
    parts := strings.Split(base, ".")
    config = parts[0]
  } else {
    config = base
  }

  return path, config
}

func (c *Config) Get(key string) interface{} {
  return c.v.Get(key)
}
func (c *Config) GetBool(key string) bool {
  return c.v.GetBool(key)
}
func (c *Config) GetFloat64(key string) float64 {
  return c.v.GetFloat64(key)
}
func (c *Config) GetInt(key string) int {
  return c.v.GetInt(key)
}
func (c *Config) GetString(key string) string {
  return c.v.GetString(key)
}
func (c *Config) GetStringMap(key string) map[string]interface{} {
  return c.v.GetStringMap(key)
}
func (c *Config) GetStringMapString(key string) map[string]string {
  return c.v.GetStringMapString(key)
}
func (c *Config) GetStringSlice(key string) []string {
  return c.v.GetStringSlice(key)
}
func (c *Config) GetTime(key string) time.Time {
  return c.v.GetTime(key)
}
func (c *Config) GetDuration(key string) time.Duration {
  return c.v.GetDuration(key)
}

func (c *Config) IsSet(key string) bool {
  return c.v.IsSet(key)
}

func (c *Config) GetOrDefault(key string, def interface{} ) interface{} {
  if i := c.v.Get(key); i != nil {
    return i
  } else {
    return def
  }
}

func (c *Config) GetBoolOrDefault(key string, def bool ) bool {
  if i := c.v.Get(key); i != nil {
    return cast.ToBool(i)
  } else {
    return def
  }
}

func (c *Config) GetFloat64OrDefault(key string, def float64 ) float64 {
  if i := c.v.Get(key); i != nil {
    return cast.ToFloat64(i)
  } else {
    return def
  }
}

func (c *Config) GetIntOrDefault(key string, def int ) int {
  if i := c.v.Get(key); i != nil {
    return cast.ToInt(i)
  } else {
    return def
  }
}

func (c *Config) GetStringOrDefault(key string, def string ) string {
  if i := c.v.Get(key); i != nil {
    return cast.ToString(i)
  } else {
    return def
  }
}

func (c *Config) GetStringMapOrDefault(key string, def map[string]interface{} ) map[string]interface{} {
  if i := c.v.Get(key); i != nil {
    return cast.ToStringMap(i)
  } else {
    return def
  }
}

func (c *Config) GetStringMapStringOrDefault(key string, def map[string]string ) map[string]string {
  if i := c.v.Get(key); i != nil {
    return cast.ToStringMapString(i)
  } else {
    return def
  }
}

func (c *Config) GetStringSliceOrDefault(key string, def []string ) []string {
  if i := c.v.Get(key); i != nil {
    return cast.ToStringSlice(i)
  } else {
    return def
  }
}

func (c *Config) GetTimeOrDefault(key string, def time.Time ) time.Time {
  if i := c.v.Get(key); i != nil {
    return cast.ToTime(i)
  } else {
    return def
  }
}

func (c *Config) GetDurationOrDefault(key string, def time.Duration ) time.Duration {
  if i := c.v.Get(key); i != nil {
    return cast.ToDuration(i)
  } else {
    return def
  }
}

// If you need to do anything simpler than loading a config, and getting
// values, then you wil likely want to work directly with the Viper instance.
func (c *Config) Viper() *viper.Viper {
  return c.v
}
