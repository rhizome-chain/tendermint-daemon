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
	flagDelMemberThreshold = "daemon_del_member_threshold"
	flagApiAddr        = "daemon_api_addr"
)

type DaemonConfig struct {
	ChainID                  string
	NodeID                   string
	NodeName                 string
	APIAddr                  string `mapstructure:"daemon_api_addr"`
	AliveThresholdBlocks     uint   `mapstructure:"daemon_alive_threshold"`
	DelMemberThresholdBlocks uint   `mapstructure:"daemon_del_member_threshold"`
}

func NewDaemonConfig() *DaemonConfig {
	conf := &DaemonConfig{
		AliveThresholdBlocks: uint(2),
		DelMemberThresholdBlocks: uint(15),
	}
	return conf
}

func AddDaemonFlags(cmd *cobra.Command) {
	cmd.Flags().Uint(flagAliveThreshold, 2, "Alive Threshold Blocks")
	cmd.Flags().Uint(flagDelMemberThreshold, 15, "Delete Member Threshold Blocks")
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
daemon_alive_threshold = "{{ .AliveThresholdBlocks }}"
daemon_del_member_threshold = "{{ .DelMemberThresholdBlocks }}"
`
