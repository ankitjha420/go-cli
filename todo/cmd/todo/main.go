package main

import (
	"flag"
	"fmt"
	"github.com/ankitjha420/todo"
	"os"
)

const todoFileName = ".todo.json"

func main() {
	// available CLI args
	task := flag.String(
		"task",
		"",
		"Task to be included in the Todo list",
	)
	list := flag.Bool("list", false, "List all tasks")
	complete := flag.Int("complete", 0, "Item to be completed")
	flag.Parse()

	l := &todo.List{}

	if err := l.Get(todoFileName); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	switch {
	case *list:
		for _, item := range *l {
			if !item.Done {
				fmt.Print(list)
			}
		}

	case *complete > 0:
		if err := l.Complete(*complete); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if err := l.Save(todoFileName); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	case *task != "":
		l.Add(*task)
		if err := l.Save(todoFileName); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	default:
		// invalid flag(s)
		_, _ = fmt.Fprintln(os.Stderr, "invalid options")
		os.Exit(1)
	}
}
