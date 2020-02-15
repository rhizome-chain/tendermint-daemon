package hello

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
	
	"github.com/rhizome-chain/tendermint-daemon/daemon/worker"
)

const (
	FactoryName = "hello-worker"
)

type Factory struct {
}

var _ worker.Factory = (*Factory)(nil)

func (fac *Factory) Name() string { return FactoryName }
func (fac *Factory) Space() string { return worker.SpaceDefaultWorker }

func (fac *Factory) NewWorker(helper *worker.Helper) (worker.Worker, error) {
	worker := &Worker{id: helper.ID(), helper: helper}
	
	infoBytes := worker.helper.Job().Data
	
	jobInfo := JobInfo{}
	err := json.Unmarshal(infoBytes, &jobInfo)
	
	if err != nil {
		err = errors.New("cannot create hello worker(unmarshal json) " + err.Error())
		helper.Error("[ERROR] Create Hello Worker(unmarshal json)", "data", string(infoBytes))
		return nil, err
	}
	worker.jobInfo = jobInfo
	
	interval, err := time.ParseDuration(worker.jobInfo.Interval)
	if err != nil {
		err = errors.New("cannot create hello worker(ParseDuration) " + err.Error())
		return nil, err
	}
	
	jobInfo.intervalDur = interval
	worker.jobInfo = jobInfo
	
	return worker, nil
}

type JobInfo struct {
	Interval    string `json:"interval"`
	Greet       string `json:"greet"`
	intervalDur time.Duration
}

type Worker struct {
	id      string
	started bool
	helper  *worker.Helper
	jobInfo JobInfo
}

var _ worker.Worker = (*Worker)(nil)

func (worker *Worker) ID() string      { return worker.id }
func (worker *Worker) IsStarted() bool { return worker.started }

func (worker *Worker) Start() error {
	fmt.Println(" - Worker Started ", worker.id)
	worker.started = true
	
	for worker.started {
		fmt.Printf(" - Worker[%s] worked. %s \n", worker.id, worker.jobInfo.Greet)
		time.Sleep(worker.jobInfo.intervalDur)
	}
	
	fmt.Println(" - Worker Stopped ", worker.id)
	return nil
}

func (worker *Worker) Stop() error {
	worker.started = false
	return nil
}
