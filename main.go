package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

type TaskExecutor struct {
	wg sync.WaitGroup

	taskCh chan string
}

var (
	nproc    int
	filePath string
)

func init() {
	flag.Parse()

	nprocArg := flag.Arg(0)
	if nprocValue, err := strconv.Atoi(nprocArg); err != nil {
		log.Fatalf("Error converting %s to int", nprocArg)
	} else if nprocValue <= 0 {
		log.Fatalf("Error NPROC to be great zero, current value: %d", nprocValue)
	} else {
		nproc = nprocValue
	}

	filePath = flag.Arg(1)
	if filePath == "" {
		log.Fatal("Error FILE does not empty")
	}
}

func main() {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	var (
		reader   = bufio.NewReader(file)
		scanner  = bufio.NewScanner(reader)
		start    = time.Now()
		executor TaskExecutor
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Создаем канал для отправки задач на обработку
	executor.taskCh = make(chan string, nproc)
	defer close(executor.taskCh)

	// Создаем обработчики выполняющие работу
	for i := 0; i < nproc; i++ {
		executor.wg.Add(1)
		go worker(ctx, &executor)
	}

	// Отправляем задачи на обработку
	for scanner.Scan() {
		executor.taskCh <- scanner.Text()
	}
	cancel()

	// Дожидаемся выполнения всех задач
	executor.wg.Wait()

	end := time.Now()
	diff := end.Sub(start)

	fmt.Printf("Duration of process: %v \n", diff)
}

func worker(ctx context.Context, executor *TaskExecutor) {
	defer executor.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return

		case duration := <-executor.taskCh:
			value, err := strconv.ParseInt(duration, 10, 16)
			if err != nil || value < 0 {
				log.Fatalf("Error converting %s to int", duration)
			}

			time.Sleep(time.Millisecond * time.Duration(value))
			fmt.Println(value)
		}
	}
}
