package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/NiksonGo/stream-processing/pipeline"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: md5 <directory>")
	}
	root := os.Args[1]

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	walker := NewFileWalker(root)
	md5 := NewMD5Node(10)
	printer := NewPrinterNode()

	md5.SetInput(walker.Output())
	printer.SetInput(md5.Output())

	graph := pipeline.NewGraph()

	pipeline.AddNode(graph, walker)
	pipeline.AddBidiNode(graph, md5)
	pipeline.AddNode(graph, printer)

	if err := graph.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
