package config

import (
  "log"
  //"github.com/spf13/pflag"
)

type ConfigFlag struct {
  Name string
  Type string
  Default interface{}
  Help string
}

type ConfigFlags []ConfigFlag


func (c *ConfigFlags) Bind() {
  //flagSet := pflag.NewFlagSet("config", pflag.ContinueOnError)

  for _, flag := range *c {
    log.Printf("Flags: %#v", flag)
  }
}
