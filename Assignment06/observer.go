package main

import "sync"

type EventType string

type TypedEvent interface{ Type() EventType }

type Observer[E TypedEvent] interface{ OnEvent(E) }

type Subject[E TypedEvent] interface {
	Subscribe(EventType, Observer[E]) func()
	Publish(E)
}

type listener[E TypedEvent] struct {
	id       int
	observer Observer[E]
}

type Bus[E TypedEvent] struct {
	mu        sync.RWMutex
	nextID    int
	listeners map[EventType][]listener[E]
}

func NewBus[E TypedEvent]() *Bus[E] {
	return &Bus[E]{listeners: make(map[EventType][]listener[E])}
}

func (b *Bus[E]) Subscribe(kind EventType, observer Observer[E]) func() {
	b.mu.Lock()
	b.nextID++
	id := b.nextID
	b.listeners[kind] = append(b.listeners[kind], listener[E]{id, observer})
	b.mu.Unlock()
	return func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		for index, item := range b.listeners[kind] {
			if item.id == id {
				b.listeners[kind] = append(b.listeners[kind][:index], b.listeners[kind][index+1:]...)
				return
			}
		}
	}
}

func (b *Bus[E]) Publish(event E) {
	b.mu.RLock()
	listeners := append([]listener[E](nil), b.listeners[event.Type()]...)
	b.mu.RUnlock()
	for _, listener := range listeners {
		listener.observer.OnEvent(event)
	}
}

type ObserverFunc[E TypedEvent] func(E)

func (f ObserverFunc[E]) OnEvent(event E) { f(event) }
