package conveyer

import (
	"context"
	"sync"
)

type Conveyer struct {
	mu       sync.RWMutex
	channels map[string]chan string
	size     int
	handlers []func(context.Context) error
}

func New(size int) *Conveyer {
	return &Conveyer{
		mu:       sync.RWMutex{},
		channels: make(map[string]chan string),
		size:     size,
		handlers: make([]func(context.Context) error, 0),
	}
}

func (conv *Conveyer) getOrCreateChannel(name string) chan string {
	conv.mu.Lock()
	defer conv.mu.Unlock()

	if ch, exists := conv.channels[name]; exists {
		return ch
	}

	ch := make(chan string, conv.size)
	conv.channels[name] = ch

	return ch
}

func (conv *Conveyer) getOrCreateChannels(names ...string) {
	for _, name := range names {
		conv.getOrCreateChannel(name)
	}
}
