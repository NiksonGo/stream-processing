package pipeline

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"sync"
)

type MD5Result struct {
	Path string
	Hash string
	Err  error
}

type MD5Node struct {
	input       <-chan string
	output      chan MD5Result
	parallelism int
}

func NewMD5Node(parallelism int, buf int) *MD5Node {
	return &MD5Node{
		output:      make(chan MD5Result, buf),
		parallelism: parallelism,
	}
}

func (n *MD5Node) SetInput(ch <-chan string) {
	fmt.Println("[MD5Node] SetInput called")
	n.input = ch
}

func (n *MD5Node) Output() <-chan MD5Result {
	return n.output
}

func (n *MD5Node) Inputs() []<-chan string {
	if n.input == nil {
		return nil
	}
	return []<-chan string{n.input}
}

func (n *MD5Node) Outputs() []chan<- MD5Result {
	return []chan<- MD5Result{n.output}
}

func (n *MD5Node) Name() string {
	return "MD5Node"
}

func (n *MD5Node) Run(ctx context.Context) error {
	fmt.Println("[MD5Node] Run started")
	defer fmt.Println("[MD5Node] Run finished")
	defer close(n.output)

	var wg sync.WaitGroup
	sem := make(chan struct{}, n.parallelism)

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case path, ok := <-n.input:
			if !ok {
				break loop
			}

			wg.Add(1)
			sem <- struct{}{}

			go func(p string) {
				defer wg.Done()
				defer func() { <-sem }()

				file, err := os.Open(p)
				if err != nil {
					select {
					case n.output <- MD5Result{Path: p, Err: err}:
					case <-ctx.Done():
					}
					return
				}
				defer file.Close()

				hash := md5.New()
				if _, err := io.Copy(hash, file); err != nil {
					select {
					case n.output <- MD5Result{Path: p, Err: err}:
					case <-ctx.Done():
					}
					return
				}

				select {
				case n.output <- MD5Result{
					Path: p,
					Hash: hex.EncodeToString(hash.Sum(nil)),
				}:
				case <-ctx.Done():
				}
			}(path)
		}
	}

	wg.Wait() //Ждём все запущенные горутины
	return ctx.Err()
}






