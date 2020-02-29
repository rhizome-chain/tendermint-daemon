package hello

import (
	"encoding/json"
	"errors"
	"github.com/rhizomata-io/dist-daemon-tendermint/types"
	tdtypes "github.com/rhizome-chain/tendermint-daemon/types"
	"time"
	
	"github.com/rhizome-chain/tendermint-daemon/daemon/worker"
)

const (
	ProxyFactoryName = "hello-proxy-worker"
)

// JobInfo
type ProxyJobInfo struct {
	Source string `json:"source"`
}

// Event
type ProxyEvent struct {
	Source string
	Log    string
	Time   time.Time
}

func init() {
	tdtypes.BasicCdc.RegisterConcrete(ProxyEvent{}, "hello-proxy/Event", nil)
}

type ProxyFactory struct {
}

var _ worker.Factory = (*ProxyFactory)(nil)

func (fac *ProxyFactory) Name() string  { return ProxyFactoryName }
func (fac *ProxyFactory) Space() string { return SpaceHello }

func (fac *ProxyFactory) NewWorker(helper *worker.Helper) (worker.Worker, error) {
	worker := &ProxyWorker{id: helper.ID(), helper: helper}
	
	infoBytes := worker.helper.Job().Data
	
	jobInfo := ProxyJobInfo{}
	err := json.Unmarshal(infoBytes, &jobInfo)
	
	if err != nil {
		err = errors.New("cannot create hello proxy worker(unmarshal json) " + err.Error())
		helper.Error("[ERROR] Create Hello Proxy Worker(unmarshal json)", "data", string(infoBytes))
		return nil, err
	}
	worker.jobInfo = jobInfo
	
	return worker, nil
}

// Worker
type ProxyWorker struct {
	id      string
	started bool
	helper  *worker.Helper
	jobInfo ProxyJobInfo
	proxy   worker.Proxy
	wait    chan bool
}

var _ worker.Worker = (*Worker)(nil)

func (worker *ProxyWorker) ID() string      { return worker.id }
func (worker *ProxyWorker) IsStarted() bool { return worker.started }

func (worker *ProxyWorker) Start() error {
	proxy, err := worker.helper.NewWorkerProxy(worker.jobInfo.Source)
	
	if err != nil {
		err = errors.New("cannot create hello proxy worker " + err.Error())
		worker.helper.Error("[ERROR] Create Hello Proxy Worker", err)
		return err
	}
	
	worker.proxy = proxy
	
	worker.helper.Info(" - Worker Started " + worker.id)
	worker.started = true
	
	worker.wait = make(chan bool)
	
	var checkpoint string
	worker.helper.GetCheckpoint(&checkpoint)
	
	//fmt.Printf("\n -- Hello Proxy[%s] checkpoint=%s\n\n", worker.ID(), checkpoint)
	
	cancelSubs, err := worker.proxy.CollectAndSubscribe("in", checkpoint, func(jobID string, topic string, rowID string, value []byte) bool {
		srcEvent := Event{}
		types.BasicCdc.UnmarshalBinaryBare(value, &srcEvent)
		bts, err := json.Marshal(srcEvent)
		if err != nil {
			worker.helper.Error("Hello ProWorker Started " + worker.id)
		} else {
			worker.doWork(rowID, string(bts))
		}
		
		return worker.started
	})
	
	if err != nil {
		worker.helper.Error("[HelloProxy] cannot subscribe to " + worker.jobInfo.Source)
		return err
	}
	
	defer cancelSubs()
	
	<-worker.wait
	
	worker.helper.Info(" - Worker Stopped " + worker.id)
	return nil
}

func (worker *ProxyWorker) doWork(rowID string, log string) bool {
	event := ProxyEvent{
		Source: worker.jobInfo.Source,
		Log:    log,
		Time:   time.Now(),
	}
	
	worker.helper.PutObject("in", rowID, event)
	worker.helper.PutCheckpoint(rowID)
	return worker.started
}

func (worker *ProxyWorker) Stop() error {
	if worker.started {
		worker.wait <- false
	}
	worker.started = false
	return nil
}
