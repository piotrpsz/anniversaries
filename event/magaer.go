package events

import (
	"fmt"
	"sync"

	"calendar/shared"
)

type EventManager struct {
	mtx       sync.Mutex
	consumers []*consumer
}

var (
	instance *EventManager
	once     sync.Once
)

func Instance() *EventManager {
	once.Do(func() {
		instance = newEventManager()
	})
	return instance
}

func newEventManager() *EventManager {
	return new(EventManager)
}

func (em *EventManager) Send(e *Event) {
	go func() {
		em.mtx.Lock()
		defer em.mtx.Unlock()

		for _, consumer := range em.consumers {
			if consumer.id() == e.Id() {
				consumer.stream() <- e
			}
		}
	}()
}

// Register rejestacja użytkownika i eventów, na które oczekuje.
// Rejestrowane są tylko eventy, które nie są jeszcze zarejestrowane.
func (em *EventManager) Register(user string, eventIDs ...EventID) <-chan *Event {
	em.mtx.Lock()
	defer em.mtx.Unlock()

	hash := shared.HashOfString(user)
	out := em.channelForConsumer(hash)
	if out == nil {
		out = make(chan *Event)
	}

	for _, eventID := range eventIDs {
		if em.isEventRegistered(hash, eventID) == -1 {
			em.consumers = append(em.consumers, newConsumer(hash, eventID, out))
		}
	}

	return out
}

// Release usunięcie z rejestru eventów dla wskazanego użytkownika.
// Usunięcie ostatniego wpisu powoduje zamknięcie kanału.
func (em *EventManager) Release(user string, eventIDs ...EventID) {
	em.mtx.Lock()
	defer em.mtx.Unlock()

	hash := shared.HashOfString(user)

	if stream := em.channelForConsumer(hash); stream != nil {
		for _, eventID := range eventIDs {
			if idx := em.isEventRegistered(hash, eventID); idx != -1 {
				em.removeConsumer(idx)
			}
		}
		if em.channelForConsumer(hash) == nil {
			fmt.Println("close channel")
			close(stream)
		}
	}
}

// channelForConsumer dla każdego konsumenta jest jeden unikalny kanał,
// jeśli jest dodawany nowy event to powinien być wysyłany na ten sam kanał.
// Jeśli taki kanał jeszcze nie istnieje zwracany jest nil.
func (em *EventManager) channelForConsumer(hash uint32) chan *Event {
	for _, consumer := range em.consumers {
		if consumer.hash() == hash {
			return consumer.stream()
		}
	}
	return nil
}

func (em *EventManager) isEventRegistered(hash uint32, id EventID) int {
	for idx, consumer := range em.consumers {
		if consumer.hash() == hash && consumer.id() == id {
			return idx
		}
	}
	return -1
}

func (em *EventManager) removeConsumer(i int) {
	data := em.consumers
	copy(data[i:], data[i+1:])
	em.consumers = data[:len(data)-1]
}
