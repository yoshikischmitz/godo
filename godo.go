package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strconv"
	"time"
)

type TasksRoot struct {
	Tasks []Task
}

type Task struct {
	Priority int64
	Content  string
	Date     time.Time
	Done     bool
	Index    int
	SubTasks []*Task
}

var (
	json_path string
	task_root TasksRoot
)

// Returns an array of Tasks, with indices
func TaskList() TasksRoot {
	file, err := ioutil.ReadFile(json_path)

	// if err is not nil file does not exist
	if err != nil {
		log.Fatal(err)
	}
	tasks_root := TasksRoot{}
	json.Unmarshal(file, &tasks_root)
	return tasks_root
}

// A struct function for Task structs. Converts a the referenced Task to a tab delimited
// String
// Example:
// task := Task {0, "get groceries",time.Now(),false}
// tas.String() //= [19]    [2014-3-16]    get groceries
func (t *Task) Print(i int) string {
	year, month, day := t.Date.Date()
	return fmt.Sprintf("[%d]\t[%d-%d-%d]\t%s", i, year, month, day, t.Content)
}

// Build a Task with Task.Content from string, with default values
func buildTask(s string) Task {
	tasks_root := TaskList()
	task := Task{
		Priority: 0,
		Content:  s,
		Date:     time.Now(),
		Done:     false,
		Index:    len(tasks_root.Tasks),
	}
	return task
}

// Add a new task to tasks file
func AddTask(task Task) error {
	task_list := TaskList()

	task_list.Tasks = append(task_list.Tasks, task)

	json, _ := json.MarshalIndent(task_list, "", "  ")

	os.Remove(json_path)
	ioutil.WriteFile(json_path, json, 0600)

	return nil
}

func AddSubTask(task *Task, index int) {
	for i := range task_root.Tasks {
		t := &task_root.Tasks[i]
		if t.Index == index {
			t.SubTasks = append(t.SubTasks, task)
		}
	}
	os.Remove(json_path)
	json, _ := json.MarshalIndent(task_root, "", "  ")
	ioutil.WriteFile(json_path, json, 0600)
}

// Recursively Print Task + SubTasks
func PrintTask(t *Task, idx int, tIdx int) {
	// print our tabs out
	for i := 0; i < tIdx; i++ {
		if tIdx != 0 {
			fmt.Print("|--")
		}
	}
	year, month, day := t.Date.Date()
	fmt.Printf("[%d]\t[%d-%d-%d]\t%s\n", idx, year, month, day, t.Content)
	for _, task := range t.SubTasks {
		tIdx += 1
		PrintTask(task, 1, tIdx)
	}
}

// Print all tasks in tasks.json
func PrintAllTasks() {
	var j int
	for i := range task_root.Tasks {
		t := task_root.Tasks[i]
		if t.Done == false {
			j += 1
			PrintTask(&t, j, 0)
		}
	}
}

// Sets the complete field to true at relative index for uncomplete tasks
func CompleteTask(index int) {
	// Return index to 0-index
	index -= 1
	task_root.Tasks[index].Done = true
	var j int
	for i := range task_root.Tasks {
		if task_root.Tasks[i].Done == false {
			if j == index {
				task_root.Tasks[i].Done = true
			}
			j += 1
		}
	}
	json, _ := json.MarshalIndent(task_root, "", "  ")

	os.Remove(json_path)
	ioutil.WriteFile(json_path, json, 0600)

	fmt.Printf("Task Marked as complete: %s\n", task_root.Tasks[index].Content)
}

func init() {
	// Assume that if we're not on Windows, we're on a *nix-like system
	// Should add more robust OS support in the future
	if runtime.GOOS == "windows" {
		json_path = os.Getenv("UserProfile") + "/My Documents/tasks.json"
	} else {
		json_path = os.Getenv("HOME") + "/tasks.json"
	}
	task_root = TaskList()
}

func main() {
	app := cli.NewApp()
	app.Name = "todo"
	app.Usage = "add, track, and complete todos from the commandline"
	app.Commands = []cli.Command{
		{
			Name:  "add",
			Usage: "add a task",
			Action: func(c *cli.Context) {
				task := buildTask(c.Args().First())
				AddTask(task)
				fmt.Printf("Task is added: %s\n", task.Content)
			},
		},
		{
			Name:  "subadd",
			Usage: "add a sub-task",
			Action: func(c *cli.Context) {
				task_id, _ := strconv.ParseInt(c.Args()[0], 10, 0)
				index := int(task_id)
				task := buildTask(c.Args()[1])
				AddSubTask(&task, index)
				fmt.Printf("Subtask added: %s\n", task.Content)
			},
		},
		{
			Name:  "ls",
			Usage: "lists all tasks",
			Action: func(c *cli.Context) {
				PrintAllTasks()
			},
		},
		{
			Name:  "complete",
			Usage: "Marks a task specified by the integer argument as complete",
			Action: func(c *cli.Context) {
				index, _ := strconv.ParseInt(c.Args().First(), 10, 0)
				CompleteTask(int(index))
			},
		},
	}
	app.Run(os.Args)
}
