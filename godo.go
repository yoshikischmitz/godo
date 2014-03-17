package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"os"
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
	json_path string    = os.Getenv("HOME") + "/tasks.json"
	task_root TasksRoot = TaskList()
)

// Takes a json string and converts it to a Task struct,(without an index)
func ParseTask(j string) Task {
	var task Task
	b := []byte(j)
	json.Unmarshal(b, &task)
	return task
}

// Returns an array of Tasks, with indices
func TaskList() TasksRoot {
	file, _ := ioutil.ReadFile(json_path)
	tasks_root := TasksRoot{}
	json.Unmarshal(file, &tasks_root)
	return tasks_root
}

// A struct function for Task structs. Converts a the referenced Task to a tab delimited
// String
// Example:
// task := Task {0, "get groceries",time.Now(),false}
// tas.String() //= [19]    [2014-3-16]    get groceries
func (t *Task) String() string {
	year, month, day := t.Date.Date()
	return fmt.Sprintf("[%d]\t[%d-%d-%d]\t%s", t.Index, year, month, day, t.Content)
}

// Build a Task with Task.Content from string, with default values
func buildTask(s string) Task {
	tasks_root := TaskList()
	task := Task{
		Priority: 0,
		Content:  s,
		Date:     time.Now(),
		Done:     false,
		Index:    len(tasks_root.Tasks) - 1,
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
func PrintTask(t *Task, idx int) {
	// print our tabs out
	for i := 0; i < idx; i++ {
		fmt.Print("\t")
	}
	fmt.Println(t)
	for _, task := range t.SubTasks {
		PrintTask(task, idx+1)
	}
}

// Print all tasks in tasks.json
func PrintAllTasks() {
	for i := range task_root.Tasks {
		t := task_root.Tasks[i]
		if t.Done == false {
			PrintTask(&t, 0)
		}
	}
}

// Marks a task as complete by setting the complete field to true in the JSON file
func CompleteTask(index int) {
	task_root.Tasks[index].Done = true

	json, _ := json.MarshalIndent(task_root, "", "  ")

	os.Remove(json_path)
	ioutil.WriteFile(json_path, json, 0600)

	fmt.Printf("Task Marked as complete: %s\n", task_root.Tasks[index].Content)
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
