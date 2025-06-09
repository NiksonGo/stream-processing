package pipeline

import (
	"context"
	"sync"
)

// Graph управляет запуском всех нод
type Graph struct {
	runFuncs []func(context.Context) error
}

// NewGraph создает новый пайплайн
func NewGraph() *Graph {
	return &Graph{}
}

// Add добавляет ноду любого типа
func AddNode[T any](g *Graph, n Node[T]){
	g.runFuncs = append(g.runFuncs, n.Run)
}

func AddBidiNode[In any, Out any](g *Graph, n BidiNode[In, Out]) {
	g.runFuncs = append(g.runFuncs, n.Run)
}

// Run запускает все ноды
func (g *Graph) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	errs := make(chan error, len(g.runFuncs))

	for _, run := range g.runFuncs {
		wg.Add(1)
		go func(fn func(context.Context) error) {
			defer wg.Done()
			if err := fn(ctx); err != nil {
				errs <- err
			}
		}(run)
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}
