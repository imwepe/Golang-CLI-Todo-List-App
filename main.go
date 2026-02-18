package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	Completed bool      `json:"completed"`
}

func humanizeTime(t time.Time) string {
	diff := time.Since(t)

	seconds := int(diff.Seconds())
	minutes := int(diff.Minutes())
	hours := int(diff.Hours())
	days := hours / 24

	switch {
	case seconds < 60:
		if seconds <= 1 {
			return "just now"
		}
		return fmt.Sprintf("%d seconds ago", seconds)

	case minutes < 60:
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)

	case hours < 24:
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)

	default:
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	}
}

func loadTasks() ([]Task, error) {
	if _, err := os.Stat("tasks.json"); os.IsNotExist(err) {
		return []Task{}, nil
	}

	data, err := os.ReadFile("tasks.json")
	if err != nil {
		return nil, err
	}

	var tasks []Task
	err = json.Unmarshal(data, &tasks)
	return tasks, err
}

func saveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile("tasks.json", data, 0644)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: tasks [add|list|complete]")
		return
	}

	command := os.Args[1]

	switch command {
	case "add":
		tasks, _ := loadTasks()

		task := Task{
			ID:        len(tasks) + 1,
			Title:     os.Args[2],
			CreatedAt: time.Now(),
			Completed: false,
		}

		tasks = append(tasks, task)
		saveTasks(tasks)

		fmt.Println("Added:", task.Title)

	case "list":
		tasks, _ := loadTasks()

		if len(tasks) == 0 {
			fmt.Println("No tasks found.")
			return
		}

		fmt.Printf("%-5s %-50s %-8s %s\n", "ID", "Task", "Status", "Created")

		for _, task := range tasks {
			status := " "
			if task.Completed {
				status = "✓"
			}

			fmt.Printf(
				"%-5d %-50s %-8s %s\n",
				task.ID,
				task.Title,
				fmt.Sprintf("[%s]", status),
				humanizeTime(task.CreatedAt),
			)
		}

	case "complete":
		if len(os.Args) < 3 {
			fmt.Println("Please provide task ID")
			return
		}

		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid task ID")
			return
		}

		tasks, _ := loadTasks()

		found := false
		for i := range tasks {
			if tasks[i].ID == id {
				tasks[i].Completed = true
				found = true
				break
			}
		}

		if !found {
			fmt.Println("Task not found")
			return
		}

		saveTasks(tasks)
		fmt.Println("Task completed:", id)

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Please provide task ID")
			return
		}

		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid task ID")
			return
		}

		tasks, _ := loadTasks()

		index := -1
		for i, task := range tasks {
			if task.ID == id {
				index = i
				break
			}
		}

		if index == -1 {
			fmt.Println("Task not found")
			return
		}

		// Delete task
		tasks = append(tasks[:index], tasks[index+1:]...)

		saveTasks(tasks)
		fmt.Println("Task deleted:", id)

	default:
		fmt.Println("Unknown command")
	}
}
