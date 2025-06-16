package pipeline

import (
	"context"
	"fmt"
)

type MapNode[T any, R any] struct {
	input  <-chan T
	output chan R
	fn     func(T) R
}

func NewMapNode[T any, R any](fn func(T) R, buf int) *MapNode[T, R] {
	return &MapNode[T, R]{
		output: make(chan R, buf),
		fn:     fn,
	}
}

func (n *MapNode[T, R]) SetInput(ch <-chan T) {
	fmt.Println("[MapNode] SetInput called")
	n.input = ch
}

func (n *MapNode[T, R]) Output() <-chan R {
	return n.output
}

func (n *MapNode[T, R]) Inputs() []<-chan T {
	if n.input == nil {
		return nil
	}
	return []<-chan T{n.input}
}

func (n *MapNode[T, R]) Outputs() []chan<- R {
	return []chan<- R{n.output}
}

func (n *MapNode[T, R]) Name() string {
	return "MapNode"
}

func (n *MapNode[T, R]) Run(ctx context.Context) error {
	fmt.Println("[MapNode] Run started")
	defer fmt.Println("[MapNode] Run finished")
	defer close(n.output)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case val, ok := <-n.input:
			if !ok {
				return nil
			}
			out := n.fn(val)
			select {
			case n.output <- out:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}




