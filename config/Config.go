// Package config is just a simple wrapper around github.com/spf13/viper
package config

import (
  "fmt"
  "log"
  "time"
  "strings"
  "path/filepath"
  "github.com/fsnotify/fsnotify"
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
func New(fileName string, appName string) *Config {

  var path string
  var config string
  if fileName == "" {
    path = "."
    config = appName
  } else {
    path, config = ParseConfig(fileName)
  }

  log.Printf("Loading config file: [%s]  from [%s]", config, path)

  v := viper.New()
  v.SetDefault("APP_NAME", appName)
  v.SetConfigName(config)

  v.AutomaticEnv()

  v.AddConfigPath(path)
  if path != "." {
    v.AddConfigPath(".")
  }
  err := v.ReadInConfig() // Find and read the config file
  if err != nil { // Handle errors reading the config file
    panic(fmt.Errorf("Fatal error config file: %s \n", err))
  }
  v.WatchConfig()
  v.OnConfigChange(func(e fsnotify.Event) {
    log.Println("Config file changed:", e.Name)
  })
  c := Config{v}
  return &c
}

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
