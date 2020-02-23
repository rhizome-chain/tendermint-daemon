package types

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"text/template"
	
	"github.com/spf13/viper"
)

type ModuleConfig interface {
	GetTemplate() string
}

type EmptyModuleConfig struct {
}

func (e EmptyModuleConfig) GetTemplate() string {
	return ``
}

var _ ModuleConfig = (*EmptyModuleConfig)(nil)

func LoadModuleConfigFile(filePath string, modCfg interface{}) (err error) {
	bts, err := ioutil.ReadFile(filePath)
	
	viper.ReadConfig(bytes.NewBuffer(bts))
	viper.Unmarshal(modCfg)
	return err
}

func WriteModuleConfigFile(filePath string, modCfg ModuleConfig) (err error) {
	var configTemplate *template.Template
	if configTemplate, err = template.New("configFileTemplate").Parse(modCfg.GetTemplate()); err != nil {
		return errors.New("create module config template:" + err.Error())
	}
	
	var buffer bytes.Buffer
	
	if err := configTemplate.Execute(&buffer, modCfg); err != nil {
		panic(err)
	}
	
	err = ioutil.WriteFile(filePath, buffer.Bytes(), 0644)
	if err != nil {
		panic(fmt.Sprintf("Write Module Config File failed: %v", err))
	}
	
	return err
}
