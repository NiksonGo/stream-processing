package main

import (
	"context"
	"fmt"
	"time"
	"github.com/NiksonGo/stream-processing/pipeline"
)

func main() {
	// Создаем контекст с тайм-аутом 30 секунд
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Создаем билдер с буферизацией каналов (размер буфера 100)
	builder := pipeline.NewBuilder().WithBufferSize(100)

	// Инициализируем ноды
	fileWalker := pipeline.NewFileWalker(".")// здесь вписать свое
	md5Node := pipeline.NewMD5Node(4, 100) // Параллелизм 4, буфер 100
	printerNode := pipeline.NewPrinterNode()

	// Настраиваем пайплайн
	builder.Source(fileWalker).
		Pipe(md5Node).
		Sink(printerNode)

	// Запускаем пайплайн
	fmt.Println("Starting pipeline...")
	if err := builder.Run(ctx); err != nil {
		fmt.Printf("Pipeline failed: %v\n", err)
	} else {
		fmt.Println("Pipeline completed successfully")
	}
}