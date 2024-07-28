package atoms

import rl "github.com/gen2brain/raylib-go/raylib"

type ChangeWindowSizeEventCallback = func(oldWindowSize rl.Vector2, newWindowSize rl.Vector2)

type EventCallback = func(...interface{})

type EventStore struct {
	callbacks map[string]map[int]EventCallback
	nextId    int
}

func NewEventStore() *EventStore {
	instance := &EventStore{
		callbacks: map[string]map[int]EventCallback{},
		nextId:    1,
	}

	return instance
}

func (evStore *EventStore) ListenToEvent(eventType string, callback EventCallback) (eventId int) {
	eventId = evStore.nextId

	if _, ok := evStore.callbacks[eventType]; ok == false {
		evStore.callbacks[eventType] = map[int]EventCallback{}
	}

	evStore.callbacks[eventType][eventId] = callback

	evStore.nextId += 1

	return
}

func (evStore *EventStore) DispatchEvent(eventType string, args ...interface{}) {
	for _, callback := range evStore.callbacks[eventType] {
		callback(args...)
	}
}

func (evStore *EventStore) RemoveRegisteredEvent(eventId int) {
	for eventType := range evStore.callbacks {
		delete(evStore.callbacks[eventType], eventId)
	}
}
