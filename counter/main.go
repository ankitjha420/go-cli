package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	lines := flag.Bool("l", false, "Count lines")
	flag.Parse()

	fmt.Println(count(os.Stdin, *lines))
}

func count(r io.Reader, lines bool) int {
	scanner := bufio.NewScanner(r)
	if !lines {
		scanner.Split(bufio.ScanWords)
	}

	var count int
	for scanner.Scan() {
		count += 1
	}

	return count
}
