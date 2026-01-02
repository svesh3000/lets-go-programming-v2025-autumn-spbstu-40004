package handlers

import (
	"context"
	"errors"
	"strings"
)

var ErrNoDecoration = errors.New("can't be decorated")

const (
	noDecoratorMarker = "no decorator"
	prefDecorated     = "decorated: "
)

func PrefixDecorator(ctx context.Context, input <-chan string, output chan<- string) error {
	defer close(output)

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
