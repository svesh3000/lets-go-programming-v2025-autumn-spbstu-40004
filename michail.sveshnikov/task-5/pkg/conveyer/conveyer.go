package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

var ErrChanNotFound = errors.New("chan not found")

const undefinedData = "undefined"

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

func (conv *Conveyer) getCreateChannel(name string) chan string {
	if ch, exists := conv.channels[name]; exists {
		return ch
	}

	ch := make(chan string, conv.size)
	conv.channels[name] = ch

	return ch
}

func (conv *Conveyer) getCreateChannels(names ...string) {
	for _, name := range names {
		conv.getCreateChannel(name)
	}
}

func (conv *Conveyer) RegisterDecorator(
	handlerFunc func(ctx context.Context, input chan string, output chan string) error,
	inName string,
	outName string,
) {
	conv.mu.Lock()
	defer conv.mu.Unlock()

	inChan := conv.getCreateChannel(inName)
	outChan := conv.getCreateChannel(outName)

	conv.handlers = append(conv.handlers, func(ctx context.Context) error {
		return handlerFunc(ctx, inChan, outChan)
	})
}

func (conv *Conveyer) RegisterMultiplexer(
	handlerFunc func(ctx context.Context, inputs []chan string, output chan string) error,
	inNames []string,
	outName string,
) {
	conv.mu.Lock()
	defer conv.mu.Unlock()

	conv.getCreateChannels(inNames...)

	inChans := make([]chan string, len(inNames))
	for i, name := range inNames {
		inChans[i] = conv.channels[name]
	}

	outChan := conv.getCreateChannel(outName)

	conv.handlers = append(conv.handlers, func(ctx context.Context) error {
		return handlerFunc(ctx, inChans, outChan)
	})
}

func (conv *Conveyer) RegisterSeparator(
	handlerFunc func(ctx context.Context, input chan string, outputs []chan string) error,
	inName string,
	outNames []string,
) {
	conv.mu.Lock()
	defer conv.mu.Unlock()

	inChan := conv.getCreateChannel(inName)
	conv.getCreateChannels(outNames...)

	outChans := make([]chan string, len(outNames))
	for i, name := range outNames {
		outChans[i] = conv.channels[name]
	}

	conv.handlers = append(conv.handlers, func(ctx context.Context) error {
		return handlerFunc(ctx, inChan, outChans)
	})
}

func (conv *Conveyer) Send(chanName string, data string) error {
	conv.mu.RLock()
	channel, exists := conv.channels[chanName]
	conv.mu.RUnlock()

	if !exists {
		return ErrChanNotFound
	}

	channel <- data

	return nil
}

func (conv *Conveyer) Recv(chanName string) (string, error) {
	conv.mu.RLock()
	channel, exists := conv.channels[chanName]
	conv.mu.RUnlock()

	if !exists {
		return "", ErrChanNotFound
	}

	data, ok := <-channel
	if !ok {
		return undefinedData, nil
	}

	return data, nil
}

func (conv *Conveyer) Run(ctx context.Context) error {
	defer func() {
		conv.mu.RLock()
		defer conv.mu.RUnlock()

		for _, ch := range conv.channels {
			close(ch)
		}
	}()

	conv.mu.RLock()

	group, groupCtx := errgroup.WithContext(ctx)
	for _, handler := range conv.handlers {
		currHandler := handler

		group.Go(func() error {
			return currHandler(groupCtx)
		})
	}

	conv.mu.RUnlock()

	if err := group.Wait(); err != nil {
		return fmt.Errorf("conveyer run failed: %w", err)
	}

	return nil
}
