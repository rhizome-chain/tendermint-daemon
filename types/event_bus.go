package types

import (
	"fmt"
	"log"
	"sync"
)

type EventBus struct {
	sync.Mutex
	scope   EventScope
	routers map[EventPath]*EventRouter
}

func newEventBus(scope EventScope) *EventBus {
	bus := &EventBus{
		Mutex:   sync.Mutex{},
		scope:   scope,
		routers: make(map[EventPath]*EventRouter),
	}
	return bus
}

func (bus *EventBus) getOrMakeRouter(path EventPath) *EventRouter {
	bus.Lock()
	router, ok := bus.routers[path]
	if !ok {
		router = NewEventRouter(path)
		bus.routers[path] = router
	}
	bus.Unlock()
	return router
}

func (bus *EventBus) Subscribe(path EventPath, name string, handler EventHandler) error {
	router := bus.getOrMakeRouter(path)
	router.Add(name, handler)
	log.Println(fmt.Sprintf("[EventBus:%s] Subscribe %s %s", bus.scope, path, name))
	return nil
}

func (bus *EventBus) Unsubscribe(path EventPath, name string) {
	router := bus.getOrMakeRouter(path)
	router.Remove(name)
	log.Println(fmt.Sprintf("[EventBus:%s] Unsubscribe %s %s", bus.scope, path, name))
}

func (bus *EventBus) Publish(event Event) {
	eventPath := event.Path()
	for path, router := range bus.routers {
		if eventPath.HasPrefix(path) {
			router.PushEvent(event)
		}
	}
}
