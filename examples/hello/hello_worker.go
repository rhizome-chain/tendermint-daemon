package hello

import (
	"encoding/json"
	"errors"
	"fmt"
	tdtypes "github.com/rhizome-chain/tendermint-daemon/types"
	"strconv"
	"time"
	
	"github.com/rhizome-chain/tendermint-daemon/daemon/worker"
)

const (
	FactoryName = "hello-worker"
)

// JobInfo
type JobInfo struct {
	Interval    string `json:"interval"`
	Greet       string `json:"greet"`
	intervalDur time.Duration
}

// Event
type Event struct {
	Greet       string
	TxIndex     int64
	Time        time.Time
}

func init() {
	tdtypes.BasicCdc.RegisterConcrete(Event{}, "hello/Event", nil)
}

type Factory struct {
}

var _ worker.Factory = (*Factory)(nil)

func (fac *Factory) Name() string  { return FactoryName }
func (fac *Factory) Space() string { return SpaceHello }

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
	
	interval, err := time.ParseDuration(jobInfo.Interval)
	if err != nil {
		err = errors.New("cannot create hello worker(ParseDuration) " + err.Error())
		return nil, err
	}
	
	jobInfo.intervalDur = interval
	worker.jobInfo = jobInfo
	
	return worker, nil
}

// Worker
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
	worker.helper.Info(" - Worker Started " + worker.id)
	worker.started = true
	
	var checkpoint string
	worker.helper.GetCheckpoint(&checkpoint)
	
	var count int
	if len(checkpoint)>0 {
		count1, err := strconv.Atoi(checkpoint)
		if err != nil {
			count = 0
		} else {
			count = count1
		}
	}
	for worker.started {
		count++
		if worker.doWork(count) {
			time.Sleep(worker.jobInfo.intervalDur)
		} else {
			break
		}
	}
	
	worker.helper.Info(" - Worker Stopped " + worker.id)
	return nil
}

func (worker *Worker) doWork(count int) bool {
	rowID := fmt.Sprintf("%07d",  count)
	event := Event{
		Greet: worker.jobInfo.Greet,
		TxIndex: int64(count),
		Time:    time.Now(),
	}
	worker.helper.PutObject("in", rowID, event)
	worker.helper.PutCheckpoint(rowID)
	return worker.started
}

func (worker *Worker) Stop() error {
	worker.started = false
	return nil
}
