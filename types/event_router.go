package types

import (
	"errors"
	"fmt"
	"github.com/rhizome-chain/tendermint-daemon/utils"
	"log"
	"sync"
)

type EventRouter struct {
	sync.Mutex
	Path       EventPath
	processors map[string]*utils.QueuedProcessor
	started    bool
}

func NewEventRouter(path EventPath) *EventRouter {
	return &EventRouter{
		Mutex:      sync.Mutex{},
		Path:       path,
		processors: make(map[string]*utils.QueuedProcessor),
		started:    false,
	}
}

func (router *EventRouter) Add(name string, handler EventHandler) error {
	router.Lock()
	processor, ok := router.processors[name]
	if ok {
		err := errors.New(fmt.Sprintf("EventRouter[%s] already has %s.", router.Path, name))
		return err
	}
	
	processor = utils.NewQueuedProcessor(name, func(event interface{}) {
		handler(event.(Event))
	})
	
	router.processors[name] = processor
	router.Unlock()
	
	log.Println(fmt.Sprintf("[EventRouter:%s] Adds %s", router.Path, name))
	
	processor.Start()
	
	return nil
}

func (router *EventRouter) Remove(name string) {
	router.Lock()
	processor, ok := router.processors[name]
	if !ok {
		err := errors.New(fmt.Sprintf("EventRouter[%s] doesn't have %s.", router.Path, name))
		log.Println(err)
		return
	}
	
	delete(router.processors, name)
	router.Unlock()
	log.Println(fmt.Sprintf("[EventRouter:%s] Remove %s", router.Path, name))
	processor.Stop()
}

func (router *EventRouter) PushEvent(event Event) {
	for _, proc := range router.processors {
		proc.Push(event)
	}
}
