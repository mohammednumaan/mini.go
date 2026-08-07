package repl

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mohammednumaan/mini.go/go-cli/internal/store"
)

func handleAdd(args []string) {
	if len(args) == 0 {
		fmt.Println("task title cannot be empty!")
		return
	}

	title := strings.Join(args, " ")
	title = strings.TrimSpace(title)

	if title == "" {
		fmt.Println("task title must not be empty!")
		return
	}

	store.AddTask(args[0])
	fmt.Println("Added task successfully!")
}

func handleToggle(args []string) {
	if len(args) == 0 {
		fmt.Println("task id cannot be empty!")
		return
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("invalid id")
		return
	}
	if err := store.ToggleTask(id); err != nil {
		fmt.Println("invalid id: ", err)
		return
	}

	fmt.Println("Toggled task successfully!")
}

func handleDelete(args []string) {
	if len(args) == 0 {
		fmt.Println("task id cannot be empty!")
		return
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("invalid id")
		return
	}
	if err := store.DeleteTask(id); err != nil {
		fmt.Println("invalid id: ", err)
		return
	}

	fmt.Println("Deleted task successfully!")

}

func handleView() {
	tasks := store.GetTasks()
	if len(tasks) == 0 {
		fmt.Println("no tasks to view!")
		return
	}

	for i, task := range tasks {
		var status string
		if task.Completed == true {
			status = "[x]"
		} else {
			status = "[ ]"
		}
		fmt.Printf("%d: %s %s\n", i+1, status, task.Title)
	}
}

func handleSave() {
	if err := store.SaveTask(); err != nil {
		fmt.Println("error saving tasks:", err)
		return
	}
	fmt.Println("saved all tasks successfully to disk!")
}

func Start() {
	scanner := bufio.NewScanner(os.Stdin)
	if err := store.LoadTasks(); err != nil {
		fmt.Println(err)
		return
	}

	for {
		fmt.Print("[todo-cli]> ")

		if !scanner.Scan() {
			break
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("error reading input:", err)
		}

		input := strings.Fields(scanner.Text())
		if len(input) == 0 {
			continue
		}

		cmd := strings.ToLower(input[0])
		args := input[1:]
		switch cmd {
		case "add":
			handleAdd(args)
		case "toggle":
			handleToggle(args)
		case "delete":
			handleDelete(args)
		case "view":
			handleView()
		case "save":
			handleSave()
		case "exit":
			handleSave()
			os.Exit(0)
			break
		default:
			fmt.Println("invalid command")
		}

	}
}
