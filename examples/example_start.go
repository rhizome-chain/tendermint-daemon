package main

import (
	cmd "github.com/rhizome-chain/tendermint-daemon/cmd/commands"
	"github.com/rhizome-chain/tendermint-daemon/daemon"
	"github.com/rhizome-chain/tendermint-daemon/examples/hello"
)

func main() {
	daemonProvider := &daemon.BaseProvider{}
	
	daemonProvider.AddModuleProvider(&hello.ModuleProvider{})
	
	cmd.DoCmd(daemonProvider)
}
