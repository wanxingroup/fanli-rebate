package config

import (
	"fmt"
	"time"
)

const rpcServiceDefaultHost = "127.0.0.1"
const rpcServiceDefaultPort = 8088
const rpcServiceDefaultConnectTimeout = time.Minute

type rpcServiceSetting struct {
	Host           string        `json:"host" yaml:"host"`
	Port           uint16        `json:"port" yaml:"port"`
	ConnectTimeout time.Duration `json:"connectTimeout" yaml:"connectTimeout"`
}

func (setting rpcServiceSetting) GetConnectTimeout() time.Duration {

	if setting.ConnectTimeout == 0 {
		return rpcServiceDefaultConnectTimeout
	}

	return setting.ConnectTimeout
}

func (setting rpcServiceSetting) GetAddress() string {

	return fmt.Sprintf("%s:%d", setting.GetHost(), setting.GetPort())
}

func (setting rpcServiceSetting) GetHost() string {

	if setting.Host == "" {
		return rpcServiceDefaultHost
	}
	return setting.Host
}

func (setting rpcServiceSetting) GetPort() uint16 {

	if setting.Port == 0 {
		return rpcServiceDefaultPort
	}

	return setting.Port
}
