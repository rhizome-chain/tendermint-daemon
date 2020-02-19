package types

import (
	"bytes"
	"errors"
	tmos "github.com/tendermint/tendermint/libs/os"
	"text/template"
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


func WriteModuleConfigFile(filePath string, modCfg ModuleConfig) (err error) {
	var configTemplate *template.Template
	if configTemplate, err = template.New("configFileTemplate").Parse(modCfg.GetTemplate()); err != nil {
		return errors.New("create module config template:" + err.Error())
	}
	
	var buffer bytes.Buffer
	modCfg.GetTemplate()
	
	if err := configTemplate.Execute(&buffer, modCfg); err != nil {
		panic(err)
	}
	
	tmos.MustWriteFile(filePath, buffer.Bytes(), 0644)
	return err
}