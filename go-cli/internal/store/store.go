package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/mohammednumaan/mini.go/go-cli/internal/types"
)

const jsonFilePath = "taskDir/tasks.json"

var tasks []*types.Task

func GetTasks() []*types.Task {
	return tasks
}

func AddTask(title string) {
	task := types.NewTask(title)
	tasks = append(tasks, task)
}

func ToggleTask(taskIdx int) error {
	if taskIdx > len(tasks) || taskIdx < 1 {
		return errors.New("invalid task idx")
	}
	idx := taskIdx - 1
	task := tasks[idx]
	task.Completed = !task.Completed
	return nil
}

func DeleteTask(taskIdx int) error {
	if taskIdx > len(tasks) || taskIdx < 1 {
		return errors.New("invalid task idx")
	}
	idx := taskIdx - 1
	tasks = append(tasks[:idx], tasks[idx+1:]...)
	return nil
}

func SaveTask() error {
	res, err := json.Marshal(tasks)
	if err != nil {
		return errors.New("failed to encode tasks into json")
	}

	dir := filepath.Dir(jsonFilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return errors.New("failed to create folder to store the json file")
	}

	if err := os.WriteFile(jsonFilePath, res, 0644); err != nil {
		return errors.New("failed to write tasks to disk")
	}

	return nil
}

func LoadTasks() error {
	data, err := os.ReadFile(jsonFilePath)
	if err != nil {
		return errors.New("failed to read tasks from disk")
	}

	if err := json.Unmarshal(data, &tasks); err != nil {
		return errors.New("failed to decode tasks from json")
	}

	if err := verifyTasks(); err != nil {
		return err
	}

	return nil
}

func verifyTasks() error {
	for _, task := range tasks {
		if task.Title == "" {
			return errors.New("corrupted json file: missing task title")
		}
	}
	return nil
}
