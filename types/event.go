package types

import (
	"bytes"
	"errors"
	"fmt"
	"log"
)

type EventScope string

const GlobalEventScope = EventScope("")

type EventPath string

func (path EventPath) HasPrefix(prefix EventPath) (ok bool) {
	return bytes.HasPrefix([]byte(path), []byte(prefix))
}

type Event interface {
	Path() EventPath
}

type EventHandler func(event Event)

type EventBusRegistry struct {
	eventBusMap map[EventScope]*EventBus
}

var (
	eventBusRegistry = newEventRegistry()
)

func newEventRegistry() *EventBusRegistry {
	reg := &EventBusRegistry{
		eventBusMap: make(map[EventScope]*EventBus),
	}
	return reg
}

func (registry *EventBusRegistry) RegisterEventBus(scope EventScope) *EventBus {
	bus, ok := registry.eventBusMap[scope]
	
	if ok {
		log.Fatalf("[EventBusRegistry] EventBus[%s] is already registered.", scope)
		return bus
	}
	
	bus = newEventBus(scope)
	
	registry.eventBusMap[scope] = bus
	
	return bus
}

func (registry *EventBusRegistry) Subscribe(scope EventScope, path EventPath, name string, handler EventHandler) error {
	bus, ok := registry.eventBusMap[scope]
	if ok {
		return bus.Subscribe(path, name, handler)
	} else {
		return errors.New(fmt.Sprintf("Unknown Event Scope %s", scope))
	}
}

func (registry *EventBusRegistry) Unsubscribe(scope EventScope, path EventPath, name string) {
	bus, ok := registry.eventBusMap[scope]
	if ok {
		bus.Unsubscribe(path, name)
	} else {
		err := errors.New(fmt.Sprintf("Unknown Event Scope %s", scope))
		log.Println("Unsubscribe ", scope, path, name, err)
	}
}

func (registry *EventBusRegistry) Publish(scope EventScope, event Event) error {
	bus, ok := registry.eventBusMap[scope]
	if ok {
		bus.Publish(event)
		return nil
	} else {
		return errors.New(fmt.Sprintf("Unknown Event Scope %s", scope))
	}
}

func RegisterEventBus(scope EventScope) *EventBus {
	return eventBusRegistry.RegisterEventBus(scope)
}

func Subscribe(scope EventScope, path EventPath, name string, handler EventHandler) error {
	return eventBusRegistry.Subscribe(scope, path, name, handler)
}

func Unsubscribe(scope EventScope, path EventPath, name string){
	eventBusRegistry.Unsubscribe(scope, path, name)
}

func Publish(scope EventScope, event Event) error {
	return eventBusRegistry.Publish(scope, event)
}
