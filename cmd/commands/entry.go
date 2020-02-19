package commands

import (
	"os"
	"path/filepath"
	
	"github.com/tendermint/tendermint/libs/cli"
	
	"github.com/rhizome-chain/tendermint-daemon/daemon"
)

const (
	DefaultBCDir = "chainroot"
)

func DoCmd(daemonProvider *daemon.BaseProvider) {
	rootCmd := InitRootCommand()
	AddInitCommand(rootCmd, daemonProvider)
	AddStartCommand(rootCmd, daemonProvider)
	
	cmd := cli.PrepareBaseCmd(rootCmd, "TM", os.ExpandEnv(filepath.Join("./", DefaultBCDir)))
	
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
