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

	walker := pipeline.NewFileWalker(root)
	md5 := pipeline.NewMD5Node(10)
	printer := pipeline.NewPrinterNode()

	md5.SetInput(walker.Output())
	printer.SetInput(md5.Output())

	graph := pipeline.NewGraph()
	pipeline.AddNode(graph, walker)         // Node[string]
    pipeline.AddBidiNode(graph, md5)        // BidiNode[string, MD5Result]
    pipeline.AddNode[pipeline.MD5Result](graph, printer)
       // Node[MD5Result]



	if err := graph.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
