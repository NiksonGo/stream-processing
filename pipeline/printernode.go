package pipeline

import (
	"context"
	"fmt"
)

type PrinterNode struct {
	input <-chan MD5Result
}

func NewPrinterNode() *PrinterNode {
	return &PrinterNode{}
}

func (p *PrinterNode) SetInput(ch <-chan MD5Result) { p.input = ch }
func (p *PrinterNode) Inputs() []<-chan MD5Result {
	if p.input == nil {
		return nil
	}
	return []<-chan MD5Result{p.input}
}
func (p *PrinterNode) Outputs() []chan<- MD5Result {
	return nil
}

func (p *PrinterNode) Name() string          { return "PrinterNode" }

func (p *PrinterNode) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case res, ok := <-p.input:
			if !ok {
				return nil
			}
			if res.Err != nil {
				fmt.Printf("❌ %s: %v\n", res.Path, res.Err)
			} else {
				fmt.Printf("✅ %s: %s\n", res.Path, res.Hash)
			}
		}
	}
}

