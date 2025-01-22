package main_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

var (
	binName  = "todo"
	fileName = ".todo.json"
)

func TestMain(m *testing.M) {
	fmt.Println("Building todo executable...")
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	build := exec.Command("go", "build", "-o", binName)
	if err := build.Run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Failed to build:", err)
		os.Exit(1)
	}

	fmt.Println("Running tests...")
	res := m.Run()

	fmt.Println("Cleaning up...")
	_ = os.Remove(binName)
	_ = os.Remove(fileName)

	os.Exit(res)
}

func TestTodoCLI(t *testing.T) {
	task := "test task number 1"
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal("Failed to get working directory:", err)
	}

	cmdPath := filepath.Join(dir, binName)

	t.Run("AddNewTask", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-task", task)
		if err := cmd.Run(); err != nil {
			t.Fatalf("Failed to add task: %v", err)
		}
	})

	t.Run("ListTasks", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-list")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to list tasks: %v", err)
		}

		expected := task + "\n"
		if string(out) != expected {
			t.Errorf("Expected %q but got %q", expected, string(out))
		}
	})
}
