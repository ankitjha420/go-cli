package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

const (
	inputFile  = "./testdata/test1.md"
	goldenFile = "./testdata/test1.md.html"
)

func TestParseContent(t *testing.T) {
	input, err := os.ReadFile(inputFile)

	if err != nil {
		t.Fatal(err)
	}

	result, err := parseContent(input, "")

	if err != nil {
		t.Fatal(err)
	}

	expected, err := os.ReadFile(goldenFile)

	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(expected, result) {
		t.Logf("golden:\n %s\n", expected)
		t.Logf("result:\n %s\n", result)
		t.Errorf("Result Content does not match golden file")
	}
}

func TestRun(t *testing.T) {
	var mockStdOut bytes.Buffer

	if err := run(inputFile, "", &mockStdOut, true); err != nil {
		t.Fatal(err)
	}

	resultFile := strings.TrimSpace(mockStdOut.String())
	result, err := os.ReadFile(resultFile)

	if err != nil {
		t.Fatal(err)
	}

	expected, err := os.ReadFile(goldenFile)

	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(expected, result) {

		t.Logf("golden:\n %s\n", expected)
		t.Logf("result:\n %s\n", result)
		t.Errorf("Result Content does not match golden file")
	}

	_ = os.Remove(resultFile)
}
func checkEqual(one, two []byte) bool {
	isWhitespace := func(b byte) bool {
		return b == ' ' || b == '\t' || b == '\n' || b == '\r' || b == '\f' || b == '\v'
	}

	var cleanOne, cleanTwo []byte

	// Clean first slice
	for _, b := range one {
		if !isWhitespace(b) {
			cleanOne = append(cleanOne, b)
		}
	}

	// Clean second slice
	for _, b := range two {
		if !isWhitespace(b) {
			cleanTwo = append(cleanTwo, b)
		}
	}

	return bytes.Equal(cleanOne, cleanTwo)
}