package pipeline

import "context"

// Node — нода с одним входом и выходом одного типа T.
// Пример: Source, Sink, Log, Filter
type Node[T any] interface {
	Inputs() []<-chan T
	Outputs() []chan<- T
	Run(ctx context.Context) error
	Name() string
}

// BidiNode — нода с входом одного типа и выходом другого.
// Пример: Map[string → int], Decode[[]byte → struct{}]
type BidiNode[In any, Out any] interface {
	Inputs() []<-chan In
	Outputs() []chan<- Out
	Run(ctx context.Context) error
	Name() string
}
