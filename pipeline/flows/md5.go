package flows

import (
	"context"

	"github.com/NiksonGo/stream-processing/pipeline"
)

// RunMD5Printer запускает поток:
// - обходит директорию
// - считает MD5-хеши файлов (параллельно)
// - печатает результат
func RunMD5Printer(ctx context.Context, root string, parallelism int) error {
	// Источник: список файлов
	walker := pipeline.NewFileWalker(root)

	// Преобразователь: MD5-хеши
	md5 := pipeline.NewMD5Node(parallelism,  64)

	// Приёмник: вывод в stdout
	printer := pipeline.NewPrinterNode()

	// Сборка пайплайна
	return pipeline.NewBuilder().
		WithBufferSize(64).              
		Source(walker).                   
		Pipe(md5).                        
		Sink(printer).                   
		Run(ctx)
}
