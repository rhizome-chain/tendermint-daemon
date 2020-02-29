package utils

import (
	"fmt"
	"log"
)

const EOP = "end"

type QueuedProcessor struct {
	Name        string
	queue       *Queue
	procHandler func(event interface{})
	running     bool
}

func NewQueuedProcessor(name string, procHandler func(event interface{})) *QueuedProcessor {
	return &QueuedProcessor{
		Name:        name,
		queue:       NewQueue(),
		procHandler: procHandler,
	}
}

func (proc *QueuedProcessor) Push(event interface{}) {
	proc.queue.Push(event)
}

func (proc *QueuedProcessor) Start() {
	proc.running = true
	go func(){
		log.Println(fmt.Sprintf("QueuedProcessor[%s] starts.",proc.Name))
		for proc.running {
			event := proc.queue.Pop()
			if event != EOP {
				proc.procHandler(event)
			}
		}
		log.Println(fmt.Sprintf("QueuedProcessor[%s] ends.",proc.Name))
	}()
}

func (proc *QueuedProcessor) Stop() {
	proc.running = false
	proc.queue.Clear()
	proc.Push(EOP)
	log.Println(fmt.Sprintf("QueuedProcessor[%s] stopped.",proc.Name))
}

func (proc *QueuedProcessor) IsRunning() bool {
	return proc.running
}

