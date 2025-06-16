package pipeline

import (
	"context"
	"fmt"
)

type FilterNode[T any] struct {
	input  <-chan T
	output chan T
	fn     func(T) bool
}

func NewFilterNode[T any](fn func(T) bool, buf int) *FilterNode[T] {
	return &FilterNode[T]{
		output: make(chan T, buf),
		fn:     fn,
	}
}

func (n *FilterNode[T]) SetInput(ch <-chan T) {
	fmt.Println("[FilterNode] SetInput called")
	n.input = ch
}

func (n *FilterNode[T]) Output() <-chan T {
	return n.output
}

func (n *FilterNode[T]) Inputs() []<-chan T {
	if n.input == nil {
		return nil
	}
	return []<-chan T{n.input}
}

func (n *FilterNode[T]) Outputs() []chan<- T {
	return []chan<- T{n.output}
}

func (n *FilterNode[T]) Name() string {
	return "FilterNode"
}

func (n *FilterNode[T]) Run(ctx context.Context) error {
	fmt.Println("[FilterNode] Run started")
	defer fmt.Println("[FilterNode] Run finished")
	defer close(n.output)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case val, ok := <-n.input:
			if !ok {
				return nil
			}
			if n.fn(val) {
				n.output <- val
			}
		}
	}
}


