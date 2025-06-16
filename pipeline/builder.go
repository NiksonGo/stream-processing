package pipeline

import (
	"context"
	"fmt"
	
)

// Builder позволяет декларативно собирать и запускать пайплайны.
type Builder struct {
	graph      *Graph
	prev       any        // последняя добавленная нода
	bufferSize int        // буфер для каналов
}

// NewBuilder создаёт новый билдер пайплайна
func NewBuilder() *Builder {
	return &Builder{
		graph:      NewGraph(),
		bufferSize: 0,
	}
}

// WithBufferSize задаёт буферизацию каналов
func (b *Builder) WithBufferSize(n int) *Builder {
	b.bufferSize = n
	return b
}

// Source — первая нода (Node[T])
func (b *Builder) Source(n any) *Builder {
	b.prev = n

	switch node := n.(type) {
	case Node[string], Node[int], Node[any], Node[MD5Result]:
		
		// Вызов AddNode с type switch
		reflectAddNode(b.graph, node)
	default:
		panic(fmt.Sprintf("Source: unsupported node type: %T", n))
	}
	return b
}

// Pipe — промежуточная нода (BidiNode[In, Out])
func (b *Builder) Pipe(n any) *Builder {
	prevOut := getOutputChannel(b.prev)
	setInputChannel(n, prevOut, b.bufferSize)

	b.prev = n

	reflectAddBidiNode(b.graph, n)

	return b
}

// Sink — финальная нода
func (b *Builder) Sink(n any) *Builder {
	prevOut := getOutputChannel(b.prev)
	setInputChannel(n, prevOut, b.bufferSize)

	
	reflectAddNode(b.graph, n)

	return b
}

func (b *Builder) Map(fn any) *Builder {
	switch fn := fn.(type) {
	case func(string) string:
		node := NewMapNode(fn, b.bufferSize)

		input := getOutputChannel(b.prev)
		node.SetInput(input.(<-chan string))

		b.graph.mu.Lock()
		defer b.graph.mu.Unlock()
		AddBidiNode(b.graph, node)
		b.prev = node
		return b

	case func(string) MD5Result:
		node := NewMapNode(fn, b.bufferSize)

		input := getOutputChannel(b.prev)
		node.SetInput(input.(<-chan string))

		b.graph.mu.Lock()
		defer b.graph.mu.Unlock()
		AddBidiNode(b.graph, node)
		b.prev = node
		return b

	// можно расширить другими fn: func(int) string и т.д.

	default:
		panic(fmt.Sprintf("Map: unsupported function type: %T", fn))
	}
}


// Run — запускает пайплайн
func (b *Builder) Run(ctx context.Context) error {
	return b.graph.Run(ctx)
}
// getOutputChannel извлекает выход из предыдущей ноды
func getOutputChannel(n any) any {
	switch node := n.(type) {
	case interface{ Output() <-chan string }:
		return node.Output()
	case interface{ Output() <-chan MD5Result }:
		return node.Output()
	default:
		panic(fmt.Sprintf("unsupported Output() type: %T", n))
	}
}

// setInputChannel подключает вход к текущей ноде
func setInputChannel(n any, ch any, buf int) {
	switch node := n.(type) {
	case interface{ SetInput(<-chan string) }:
		node.SetInput(ch.(<-chan string))
	case interface{ SetInput(<-chan MD5Result) }:
		node.SetInput(ch.(<-chan MD5Result))
	default:
		panic(fmt.Sprintf("SetInput: unsupported node type: %T", n))
	}
}

// reflectAddNode вызывает AddNode через reflect
func reflectAddNode(g *Graph, node any) {
	switch n := node.(type) {
	case Node[string]:
		AddNode(g, n)
	case Node[MD5Result]:
		AddNode(g, n)
	default:
		panic(fmt.Sprintf("reflectAddNode: unsupported node type: %T", node))
	}
}

// reflectAddBidiNode вызывает AddBidiNode через reflect
func reflectAddBidiNode(g *Graph, node any) {
	switch n := node.(type) {
	case BidiNode[string, MD5Result]:
		AddBidiNode(g, n)
	default:
		panic(fmt.Sprintf("reflectAddBidiNode: unsupported node type: %T", node))
	}
}

func (b *Builder) Filter(fn any) *Builder {
	switch fn := fn.(type) {
	case func(string) bool:
		node := NewFilterNode(fn, b.bufferSize)

		input := getOutputChannel(b.prev)
		node.SetInput(input.(<-chan string))

		b.graph.mu.Lock()
		defer b.graph.mu.Unlock()
		AddBidiNode(b.graph, node)
		b.prev = node
		return b

	case func(MD5Result) bool:
		node := NewFilterNode(fn, b.bufferSize)

		input := getOutputChannel(b.prev)
		node.SetInput(input.(<-chan MD5Result))

		b.graph.mu.Lock()
		defer b.graph.mu.Unlock()
		AddBidiNode(b.graph, node)
		b.prev = node
		return b

	default:
		panic(fmt.Sprintf("Filter: unsupported function type: %T", fn))
	}
}

