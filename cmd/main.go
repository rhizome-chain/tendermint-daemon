package main

import (
	cmd "github.com/rhizome-chain/tendermint-daemon/cmd/commands"
	"github.com/rhizome-chain/tendermint-daemon/daemon"
)

func main() {
	daemonProvider := &daemon.BaseProvider{}
	cmd.DoCmd(daemonProvider)
}
