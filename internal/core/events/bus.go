// internal/core/events/bus.go
package events

import (
	"fmt"
	"sync"
	"time"
)

type Event interface {
	Type() string
	Timestamp() time.Time
}

type BaseEvent struct {
	timestamp time.Time
}

func (e BaseEvent) Timestamp() time.Time {
	return e.timestamp
}

type EventHandler func(event Event) error

type EventBus struct {
	subscribers map[string][]EventHandler
	queue       chan Event
	mu          sync.RWMutex
	done        chan struct{}
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
		done:        make(chan struct{}),
		metrics:     &EventMetrics{},
	}

	// Запускаем обработчик очереди событий
	go bus.processEvents()

	return bus
}

func (eb *EventBus) processEvents() {
	for {
		select {
		case event := <-eb.queue:
			eb.handleEvent(event)
		case <-eb.done:
			return
		}
	}
}

func (eb *EventBus) handleEvent(event Event) {
	eb.mu.RLock()
	handlers, exists := eb.subscribers[event.Type()]
	// Копируем хендлеры для безопасного выполнения
	handlersCopy := make([]EventHandler, len(handlers))
	copy(handlersCopy, handlers)
	eb.mu.RUnlock()

	if !exists {
		return
	}

	// Выполняем все хендлеры параллельно
	var wg sync.WaitGroup
	for _, handler := range handlersCopy {
		wg.Add(1)
		go func(h EventHandler) {
			defer wg.Done()

			if err := h(event); err != nil {
				eb.metrics.mu.Lock()
				eb.metrics.Errors++
				eb.metrics.mu.Unlock()
				// Логируем ошибку, но не прерываем выполнение
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

	// Возвращаем функцию отписки
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
	close(eb.done)
	close(eb.queue)
}
