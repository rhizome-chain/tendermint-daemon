package utils

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestQueuedProcessor(t *testing.T) {
	
	proc := NewQueuedProcessor("test", func(event interface{}) {
		text := event.(string)
		fmt.Println(text)
	})
	
	proc.Start()
	
	wg := sync.WaitGroup{}
	wg.Add(4)
	
	go func(){
		push("A", proc)
		wg.Done()
	}()
	go func(){
		push("B", proc)
		wg.Done()
	}()
	go func(){
		push("C", proc)
		wg.Done()
	}()
	go func(){
		push("D", proc)
		wg.Done()
	}()
	
	wg.Wait()
}

func push(name string, proc *QueuedProcessor) {
	for i := 0; i < 10; i++ {
		proc.Push(fmt.Sprintf("%s-%d", name, i))
		time.Sleep(1 * time.Millisecond)
	}
}
