package main

import (
	cmd "github.com/rhizome-chain/tendermint-daemon/cmd/commands"
	"github.com/rhizome-chain/tendermint-daemon/daemon"
	"github.com/rhizome-chain/tendermint-daemon/daemon/worker"
	"github.com/rhizome-chain/tendermint-daemon/examples/hello"
)

func main() {
	daemonProvider := &daemon.BaseProvider{}
	daemonProvider.OnDaemonStarted = onDaemonStarted
	daemonProvider.Factories = []worker.Factory{&hello.Factory{}}
	cmd.DoCmd(daemonProvider)
}

func onDaemonStarted(dm *daemon.Daemon) {
	
	// time.Sleep(4 * time.Second)
	// fmt.Println("\n------------- Sample STARTED -------------\n")
	//
	// // dm.GetJobRepository().RemoveAllJobs()
	//
	// jobInfo1 := []byte(`{"interval":"0.2s","greet":"hello 0.2s" }`)
	// dm.GetJobRepository().PutJobIfNotExist(job.NewWithID(hello.FactoryName, "hi1", jobInfo1))
	//
	// jobInfo2 := []byte(`{"interval":"1.5s","greet":"hello 1.5s" }`)
	// dm.GetJobRepository().PutJobIfNotExist(job.NewWithID(hello.FactoryName, "hi2", jobInfo2))
	//
	// jobInfo3 := []byte(`{"interval":"100ms","greet":"hello 100ms" }`)
	// dm.GetJobRepository().PutJobIfNotExist(job.NewWithID(hello.FactoryName, "hi3", jobInfo3))
	//
	// jobInfo4 := []byte(`{"interval":"0.8s","greet":"hi 0.8s" }`)
	// dm.GetJobRepository().PutJobIfNotExist(job.NewWithID(hello.FactoryName, "hi4", jobInfo4))
	//
	// jobInfo5 := []byte(`{"interval":"150ms","greet":"hi 150ms" }`)
	// dm.GetJobRepository().PutJobIfNotExist(job.NewWithID(hello.FactoryName, "hi5", jobInfo5))
	//
	// dm.GetJobRepository().GetJob("hi2")
}
