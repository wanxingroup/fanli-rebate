package config

import (
	launcherConfig "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/launcher/config"
)

const (
	DefaultLogLevel = "info"
)

type configuration struct {
	launcherConfig.StandardConfig `mapstructure:",squash"`
	Environment                   Environment `mapstructure:"environment"`
}

func (config *configuration) GetEnvironment() Environment {

	if config.Environment != "" {
		return config.Environment
	}

	return Release
}

var (
	Config *configuration
)

func init() {
	Config = &configuration{}
}
