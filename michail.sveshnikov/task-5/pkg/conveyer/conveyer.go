package conveyer

import (
	"context"
	"errors"
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

func (conv *Conveyer) RegisterDecorator(
	fn func(ctx context.Context, input chan string, output chan string) error,
	inName string,
	outName string,
) {
	conv.mu.Lock()
	defer conv.mu.Unlock()

	conv.getOrCreateChannels(inName, outName)
	inChan := conv.channels[inName]
	outChan := conv.channels[outName]

	conv.handlers = append(conv.handlers, func(ctx context.Context) error {
		return fn(ctx, inChan, outChan)
	})
}

func (conv *Conveyer) RegisterMultiplexer(
	fn func(ctx context.Context, inputs []chan string, output chan string) error,
	inNames []string,
	outName string,
) {
	conv.mu.Lock()
	defer conv.mu.Unlock()

	conv.getOrCreateChannels(inNames...)
	conv.getOrCreateChannel(outName)

	inChans := make([]chan string, len(inNames))
	for i, name := range inNames {
		inChans[i] = conv.channels[name]
	}

	outChan := conv.channels[outName]

	conv.handlers = append(conv.handlers, func(ctx context.Context) error {
		return fn(ctx, inChans, outChan)
	})
}

func (conv *Conveyer) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string) error,
	inName string,
	outNames []string,
) {
	conv.mu.Lock()
	defer conv.mu.Unlock()

	conv.getOrCreateChannel(inName)
	conv.getOrCreateChannels(outNames...)

	inChan := conv.channels[inName]

	outChans := make([]chan string, len(outNames))
	for i, name := range outNames {
		outChans[i] = conv.channels[name]
	}

	conv.handlers = append(conv.handlers, func(ctx context.Context) error {
		return fn(ctx, inChan, outChans)
	})
}

var (
	ErrChanNotFound = errors.New("chan not found")
)

const undefinedData = "undefined"

func (conv *Conveyer) Send(chanName string, data string) error {
	conv.mu.RLock()
	ch, exists := conv.channels[chanName]
	conv.mu.RUnlock()

	if !exists {
		return ErrChanNotFound
	}

	ch <- data

	return nil
}

func (conv *Conveyer) Recv(chanName string) (string, error) {
	conv.mu.RLock()
	ch, exists := conv.channels[chanName]
	conv.mu.RUnlock()

	if !exists {
		return "", ErrChanNotFound
	}

	data, ok := <-ch
	if !ok {
		return undefinedData, nil
	}

	return data, nil
}
