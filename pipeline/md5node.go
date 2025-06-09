package pipeline

import (
	"context"
	"crypto/md5"
	"encoding/hex"
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

func NewMD5Node(parallelism int) *MD5Node {
	return &MD5Node{
		output:      make(chan MD5Result),
		parallelism: parallelism,
	}
}

func (n *MD5Node) SetInput(ch <-chan string) { n.input = ch }
func (n *MD5Node) Output() <-chan MD5Result  { return n.output }

func (n *MD5Node) Inputs() []<-chan string {
	if n.input == nil {
		return nil
	}
	return []<-chan string{n.input}
}
func (n *MD5Node) Outputs() []chan<- MD5Result { return []chan<- MD5Result{n.output} }
func (n *MD5Node) Name() string                { return "MD5Node" }

func (n *MD5Node) Run(ctx context.Context) error {
	defer close(n.output)

	var wg sync.WaitGroup
	sem := make(chan struct{}, n.parallelism)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case path, ok := <-n.input:
			if !ok {
				wg.Wait()
				return nil
			}

			wg.Add(1)
			sem <- struct{}{}

			go func(p string) {
				defer wg.Done()
				defer func() { <-sem }()

				file, err := os.Open(p)
				if err != nil {
					n.output <- MD5Result{Path: p, Err: err}
					return
				}
				defer file.Close()

				hash := md5.New()
				if _, err := io.Copy(hash, file); err != nil {
					n.output <- MD5Result{Path: p, Err: err}
					return
				}
				n.output <- MD5Result{
					Path: p,
					Hash: hex.EncodeToString(hash.Sum(nil)),
				}
			}(path)
		}
	}
}