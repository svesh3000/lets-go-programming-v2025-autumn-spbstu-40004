package handlers

import (
	"context"
	"errors"
	"strings"
)

var (
	ErrNoDecoration     = errors.New("can't be decorated")
	ErrNoOutputChannels = errors.New("no output channels provided")
)

const (
	noDecoratorMarker = "no decorator"
	prefDecorated     = "decorated: "
)

func PrefixDecorator(ctx context.Context, input <-chan string, output chan<- string) error {
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

func Separator(ctx context.Context, input <-chan string, outputs []chan<- string) error {
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
