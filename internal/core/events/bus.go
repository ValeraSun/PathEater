package events

import (
	"fmt"
	"sync"
)

type Event interface {
	Type() string
}

type EventHandler func(event Event) error

type EventBus struct {
	subscribers map[string][]EventHandler
	queue       chan Event
	mu          sync.RWMutex
	done        bool
	metrics     *EventMetrics
}

type EventMetrics struct {
	Published int64
	Processed int64
	Errors    int64
	mu        sync.Mutex
}

func NewEventBus(bufferSize int) *EventBus {
	bus := &EventBus{
		subscribers: make(map[string][]EventHandler),
		queue:       make(chan Event, bufferSize),
		metrics:     &EventMetrics{},
	}

	return bus
}

func (eb *EventBus) ProcessEvents() {
	for event := range eb.queue {
		eb.handleEvent(event)
	}
}

func (eb *EventBus) handleEvent(event Event) {
	eb.mu.RLock()
	handlers, exists := eb.subscribers[event.Type()]
	handlersCopy := make([]EventHandler, len(handlers))
	copy(handlersCopy, handlers)
	eb.mu.RUnlock()

	if !exists {
		return
	}

	var wg sync.WaitGroup
	for _, handler := range handlersCopy {
		wg.Add(1)
		go func(h EventHandler) {
			defer wg.Done()
			err := h(event)
			if err != nil {
				eb.metrics.mu.Lock()
				eb.metrics.Errors++
				eb.metrics.mu.Unlock()
			}

			eb.metrics.mu.Lock()
			eb.metrics.Processed++
			eb.metrics.mu.Unlock()
		}(handler)
	}
	wg.Wait()
}

func (eb *EventBus) Subscribe(eventType string, handler EventHandler) (func(), error) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if eb.subscribers[eventType] == nil {
		eb.subscribers[eventType] = make([]EventHandler, 0)
	}

	eb.subscribers[eventType] = append(eb.subscribers[eventType], handler)

	unsubscribe := func() {
		eb.mu.Lock()
		defer eb.mu.Unlock()

		handlers := eb.subscribers[eventType]
		for i, h := range handlers {
			if fmt.Sprintf("%p", h) == fmt.Sprintf("%p", handler) {
				eb.subscribers[eventType] = append(handlers[:i], handlers[i+1:]...)
				break
			}
		}
	}

	return unsubscribe, nil
}

func (eb *EventBus) Publish(event Event) error {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	if eb.done {
		return nil
	}

	select {
	case eb.queue <- event:
		eb.metrics.mu.Lock()
		eb.metrics.Published++
		eb.metrics.mu.Unlock()
		return nil
	default:
		return fmt.Errorf("event bus queue full, dropping event: %s", event.Type())
	}
}

// PublishSync синхронно публикует событие (блокирует до обработки)
func (eb *EventBus) PublishSync(event Event) {
	eb.handleEvent(event)
}

func (eb *EventBus) Close() {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if eb.done {
		return
	}
	eb.done = true
	close(eb.queue)
}
