package pipeline

import (
	"context"
	"os"
	"path/filepath"
)

// FileWalker — читает список файлов из директории
type FileWalker struct {
	root   string
	output chan string
}

func NewFileWalker(root string) *FileWalker {
	return &FileWalker{
		root:   root,
		output: make(chan string),
	}
}

func (fw *FileWalker) Inputs() []<-chan string {
	return nil
}

func (fw *FileWalker) Outputs() []chan<- string {
	return []chan<- string{fw.output}
}

func (fw *FileWalker) Output() <-chan string {
	return fw.output
}

func (fw *FileWalker) Name() string {
	return "FileWalker"
}

func (fw *FileWalker) Run(ctx context.Context) error {
	println("[FileWalker] Run started")
	defer println("[FileWalker] Run finished")
	defer close(fw.output)

	walkErr := filepath.Walk(fw.root, func(path string, info os.FileInfo, err error) error {
		// если context уже отменён — сразу выйти
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err != nil {
			return err
		}
		if !info.IsDir() {
			select {
			case fw.output <- path:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		return nil
	})

	// если context отменён в обход Walk (например до старта)
	if walkErr == nil && ctx.Err() != nil {
		return ctx.Err()
	}

	return walkErr
}

