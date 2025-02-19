package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
)

func main() {
	// Verify and parse arguments
	op := flag.String("op", "sum", "Operation to be executed")
	column := flag.Int("col", 1, "CSV column on which to execute operation")

	flag.Parse()

	if err := run(flag.Args(), *op, *column, os.Stdout); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(filenames []string, op string, column int, out io.Writer) error {
	var opFunc statsFunc

	if len(filenames) == 0 {
		return ErrNoFiles
	}

	if column < 1 {
		return fmt.Errorf("%w: %d", ErrInvalidColumn, column)
	}

	switch op {
	case "sum":
		opFunc = sum
	case "avg":
		opFunc = avg
	case "min":
		opFunc = colMin
	case "max":
		opFunc = colMax
	default:
		return fmt.Errorf("%w: %s", ErrInvalidOperation, op)
	}

	consolidate := make([]float64, 0)
	// wait-group and channels for passing around csv values, errors and results ->
	wg := sync.WaitGroup{}
	filesCh := make(chan string)
	resCh := make(chan []float64)
	errCh := make(chan error)
	doneCh := make(chan struct{})

	go func() {
		defer close(filesCh)
		for _, fname := range filenames {
			filesCh <- fname
		}
	}()

	// one worker goroutine for each file ->
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for fname := range filesCh {
				file, err := os.Open(fname)
				if err != nil {
					errCh <- fmt.Errorf("cant open file: %w", err)
					return
				}

				data, err := csv2float(file, column)
				if err != nil {
					errCh <- err
				}

				if err = file.Close(); err != nil {
					errCh <- err
				}
				resCh <- data
			}
		}()
	}

	go func() {
		wg.Wait()
		close(doneCh)
	}()

	for {
		select {
		case err := <-errCh:
			return err
		case data := <-resCh:
			consolidate = append(consolidate, data...)
		case <-doneCh:
			_, err := fmt.Fprintln(out, opFunc(consolidate))
			return err
		}
	}
}
