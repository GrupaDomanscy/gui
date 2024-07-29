package atoms

import rl "github.com/gen2brain/raylib-go/raylib"

type ChangeWindowSizeEventCallback = func(oldWindowSize rl.Vector2, newWindowSize rl.Vector2)

type EventCallback = func(...interface{})

type EventBus struct {
	callbacks map[string]map[int]EventCallback
	nextId    int
}

func NewEventBus() *EventBus {
	instance := &EventBus{
		callbacks: map[string]map[int]EventCallback{},
		nextId:    1,
	}

	return instance
}

func (evStore *EventBus) ListenToEvent(eventType string, callback EventCallback) (eventId int) {
	eventId = evStore.nextId

	if _, ok := evStore.callbacks[eventType]; !ok {
		evStore.callbacks[eventType] = map[int]EventCallback{}
	}

	evStore.callbacks[eventType][eventId] = callback

	evStore.nextId += 1

	return
}

func (evStore *EventBus) DispatchEvent(eventType string, args ...interface{}) {
	for _, callback := range evStore.callbacks[eventType] {
		callback(args...)
	}
}

func (evStore *EventBus) RemoveRegisteredEvent(eventId int) {
	for eventType := range evStore.callbacks {
		delete(evStore.callbacks[eventType], eventId)
	}
}
