package todo_test

import (
	"github.com/ankitjha420/todo"
	"os"
	"testing"
)

func TestAdd(t *testing.T) {
	list := todo.List{}
	task := "new task"
	list.Add(task)

	if list[0].Task != task {
		t.Errorf("expected %s but got %s", task, list[0].Task)
	}
}

func TestComplete(t *testing.T) {
	list := todo.List{}
	task := "new task"
	list.Add(task)

	if list[0].Task != task {
		t.Errorf("expected %s but got %s", task, list[0].Task)
	}

	if list[0].Done {
		t.Error("expected false but got true")
	}

	err := list.Complete(1)
	if err != nil {
		t.Error(err.Error())
	} else if list[0].Done != true {
		t.Error("expected true but got false")
	}
}

func TestDelete(t *testing.T) {
	list := todo.List{}

	tasks := []string{
		"one",
		"two",
		"three",
	}
	for _, v := range tasks {
		list.Add(v)
	}

	for i, v := range list {
		if v.Task != tasks[i] {
			t.Errorf("expected %s but got %s", tasks[i], v.Task)
		}
	}

	err := list.Delete(2)
	if err != nil {
		t.Error(err.Error())
	} else if len(list) != 2 {
		t.Errorf("expected length 2 but got %d", len(list))
	}
}

func TestSaveGet(t *testing.T) {
	l1 := todo.List{}
	l2 := todo.List{}

	task := "new task"
	l1.Add(task)
	if l1[0].Task != task {
		t.Errorf("expected %s but got %s", task, l1[0].Task)
	}

	temp, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatal("could not create temporary file")
	}
	defer os.Remove(temp.Name())

	if err := l1.Save(temp.Name()); err != nil {
		t.Fatalf("error saving to temporary file: %e", err)
	}

	if err := l2.Get(temp.Name()); err != nil {
		t.Fatalf("error reading temporary file: %e", err)
	}

	if l1[0].Task != l2[0].Task {
		t.Errorf("task %s did not match %s", l1[0].Task, l2[0].Task)
	}
}
