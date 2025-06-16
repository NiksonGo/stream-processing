# Stream Processing Library (Go)

Простая и расширяемая библиотека на Go для построения **потоковых пайплайнов обработки данных** с использованием композиции узлов, декларативной сборки и контроля параллелизма.

---

## Возможности

* Декларативная сборка пайплайнов через `Builder`
* Переиспользуемые типы узлов:

  * `FileWalker` — рекурсивный обход файлов
  * `MD5Node` — параллельный расчёт MD5-хешей
  * `PrinterNode` — вывод в консоль
* Поддержка:

  * Fan-out / fan-in схем
  * Буферизированных каналов
  * Управляемого параллелизма
* Корректное завершение через `context.Context`

---

## Установка

```bash
go get github.com/NiksonGo/stream-processing
```

---

## Пример использования

### Вычисление MD5-хешей для всех файлов в директории

```go
ctx := context.Background()

err := flows.RunMD5Printer(ctx, "/home/user/data", 10)
if err != nil {
    log.Fatal(err)
}
```

Или запуск из терминала:

```bash
go run example/md5/main.go /home/user/data
```

---

## Демонстрационная задача (из технического задания)

**Задача:** Построить пайплайн, который:

* Рекурсивно обходит директорию
* Вычисляет MD5-хеши для каждого файла
* Выполняет расчёты параллельно
* Позволяет задавать степень параллелизма (по умолчанию: 10)

**Статус:** ✅ Полностью реализовано

### Схема пайплайна

```
FileWalker --> MD5Node (parallel) --> PrinterNode
```

---

## Разработка

```
example/md5/           # Пример использования
pipeline/              # Основная библиотека
pipeline/flows/        # Готовые сценарии (например, RunMD5Printer)
```

Запуск:

```bash
# Клонировать репозиторий
git clone https://github.com/NiksonGo/stream-processing.git
cd stream-processing

# Создать директорию с файлами
mkdir testdir
echo hello > testdir/a.txt
echo world > testdir/b.txt

# Запустить пайплайн с указанием директории
go run example/md5/main.go ./testdir
```

---
## Документация API

### FileWalker

* Обходит директорию рекурсивно
* Выдаёт имена всех найденных файлов

### MD5Node

* Получает пути к файлам
* Вычисляет MD5-хеши
* Работает параллельно (настраиваемо)

### PrinterNode

* Принимает `MD5Result`
* Выводит путь к файлу и его хеш в stdout

### Builder

* Позволяет декларативно собирать пайплайн:

```go
NewBuilder().
  WithBufferSize(64).
  Source(walker).
  Pipe(md5).
  Sink(printer).
  Run(ctx)
```

---




