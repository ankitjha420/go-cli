package todo

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Todo struct {
	Task        string    `json:"task"`
	Done        bool      `json:"done"`
	CreatedAt   time.Time `json:"createdAt"`
	CompletedAt time.Time `json:"completedAt"`
}

type List []Todo

func (list *List) Add(task string) {
	todo := Todo{
		Task: task,
		Done: false,
		CreatedAt: time.Now(),
		CompletedAt: time.Time{},
	}
	*list = append(*list, todo)
}

func (list *List) Complete(i int) error {
	if i <= 0 || i > len(*list) {
		return fmt.Errorf("todo item %d does not exist", i)
	}

	l := *list
	l[i - 1].Done = true
	l[i - 1].CompletedAt = time.Now()

	return nil
}

func (list *List) Delete(i int) error {
	if i <= 0 || i > len(*list) {
		return fmt.Errorf("todo item %d does not exist", i)
	}

	l := *list
	*list = append(l[:i - 1], l[i:]...)

	return nil
}

func (list *List) Save(filename string) error {
	j, err := json.MarshalIndent(list, "", "\t")
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, j, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (list *List) Get(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	if len(file) == 0 {
		return nil
	}

	return json.Unmarshal(file, list)
}

func(list *List) String() string {
	formatted := ""

	for k, t := range *list {
		prefix := "  "
		if t.Done {
			prefix = "X  "
		}

		formatted += fmt.Sprintf("%s%d:%s\n", prefix, k+1, t.Task)
	}

	return formatted
}