package pipeline

import (
	"context"
	"sync"
)

// Graph управляет запуском всех нод в пайплайне
type Graph struct {
	mu       sync.Mutex
	runFuncs []func(ctx context.Context) error
}

// NewGraph создает новый граф (пайплайн)
func NewGraph() *Graph {
	return &Graph{}
}

// AddNode добавляет обычную ноду (Node[T])
func AddNode[T any](g *Graph, n Node[T]) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.runFuncs = append(g.runFuncs, n.Run)
}

// AddBidiNode добавляет промежуточную ноду (BidiNode[In, Out])
func AddBidiNode[In any, Out any](g *Graph, n BidiNode[In, Out]) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.runFuncs = append(g.runFuncs, n.Run)
}

// Run запускает все ноды и возвращает первую ошибку или nil
func (g *Graph) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	errs := make(chan error, len(g.runFuncs))

	for _, run := range g.runFuncs {
		wg.Add(1)
		go func(fn func(context.Context) error) {
			defer wg.Done()
			if err := fn(ctx); err != nil {
				// отправляем ошибку только один раз
				select {
				case errs <- err:
					cancel() //Останавливаем всё
				default:
				}
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
