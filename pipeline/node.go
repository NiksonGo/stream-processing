package pipeline

import "context"

// Node[T] описывает потоковую ноду с входами/выходами типа T
type Node[T any] interface {
	Inputs() []<-chan T
	Outputs() []chan<- T
	Run(ctx context.Context) error
	Name() string
}


type BidiNode[In any, Out any] interface {
	Inputs() []<-chan In
	Outputs() []chan<- Out
	Run(ctx context.Context) error
	Name() string
}
