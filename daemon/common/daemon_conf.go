package common

import (
	"bytes"
	"errors"
	"fmt"
	
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	
	"io/ioutil"
	"text/template"
)

const (
	SpaceDaemon = "daemon"
)

const (
	flagAliveThreshold = "daemon_alive_threshold"
	flagApiAddr        = "daemon_api_addr"
)

type DaemonConfig struct {
	ChainID               string
	NodeID                string
	NodeName              string
	APIAddr               string `mapstructure:"daemon_api_addr"`
	AliveThresholdSeconds uint   `mapstructure:"daemon_alive_threshold"`
}

func NewDaemonConfig() *DaemonConfig {
	conf := &DaemonConfig{
		AliveThresholdSeconds: uint(2),
	}
	return conf
}

func AddDaemonFlags(cmd *cobra.Command) {
	cmd.Flags().Uint(flagAliveThreshold, 2, "Alive Threshold Seconds")
	cmd.Flags().String(flagApiAddr, "0.0.0.0:7777", "API Service ip:port")
}

func LoadConfigFile(filePath string, dmCfg *DaemonConfig) (err error) {
	bts, err := ioutil.ReadFile(filePath)
	
	viper.ReadConfig(bytes.NewBuffer(bts))
	viper.Unmarshal(dmCfg)
	
	return err
}

func WriteConfigFile(filePath string, dmCfg *DaemonConfig) (err error) {
	var configTemplate *template.Template
	if configTemplate, err = template.New("configFileTemplate").Parse(daemonConfigTemplate); err != nil {
		return errors.New("create module config template:" + err.Error())
	}
	
	var buffer bytes.Buffer
	
	if err := configTemplate.Execute(&buffer, dmCfg); err != nil {
		panic(err)
	}
	
	err = ioutil.WriteFile(filePath, buffer.Bytes(), 0644)
	if err != nil {
		panic(fmt.Sprintf("Write Module Config File failed: %v", err))
	}
	
	return err
}

var daemonConfigTemplate = `#Daemon config
# This is a TOML config file.
# For more information, see https://github.com/toml-lang/toml

daemon_api_addr = "{{ .APIAddr }}"
daemon_alive_threshold = "{{ .AliveThresholdSeconds }}"
`
