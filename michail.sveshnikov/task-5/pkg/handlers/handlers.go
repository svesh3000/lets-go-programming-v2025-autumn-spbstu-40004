package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var (
	ErrNoDecoration     = errors.New("can't be decorated")
	ErrNoOutputChannels = errors.New("no output channels provided")
)

const (
	noDecoratorMarker = "no decorator"
	prefDecorated     = "decorated: "
	noMultiplexer     = "no multiplexer"
)

func PrefixDecoratorFunc(ctx context.Context, input <-chan string, output chan<- string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, noDecoratorMarker) {
				return ErrNoDecoration
			}

			processed := data
			if !strings.HasPrefix(data, prefDecorated) {
				processed = prefDecorated + data
			}

			select {
			case <-ctx.Done():
				return nil
			case output <- processed:
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input <-chan string, outputs []chan<- string) error {
	outputsLen := len(outputs)
	if outputsLen == 0 {
		return ErrNoOutputChannels
	}

	currIdx := 0
	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			select {
			case <-ctx.Done():
				return nil
			case outputs[currIdx%outputsLen] <- data:
				currIdx++
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []<-chan string, output chan<- string) error {
	if len(inputs) == 0 {
		return nil
	}

	var wg sync.WaitGroup
	wg.Add(len(inputs))

	for i := range inputs {
		go func(inputChan <-chan string) {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case data, ok := <-inputChan:
					if !ok {
						return
					}

					if !strings.Contains(data, noMultiplexer) {
						select {
						case <-ctx.Done():
							return
						case output <- data:
						}
					}
				}
			}
		}(inputs[i])
	}

	wg.Wait()

	return nil
}
