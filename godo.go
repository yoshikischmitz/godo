package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Tasks struct {
	json_path string
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
	json_path string = os.Getenv("HOME") + "/tasks.json"
)

// Takes a json string and converts it to a Task struct,(without an index)
func ParseTask(j string) Task {
	var task Task
	b := []byte(j)
	json.Unmarshal(b, &task)
	return task
}

// Returns an array of Tasks, with indices
func TaskList() []Task {
	file, _ := ioutil.ReadFile(json_path)
	file_str := string(file)
	task_str_slice := strings.Split(file_str, "\n")
	task_list := make([]Task, len(task_str_slice)-1)
	for i := range task_str_slice {
		if task_str_slice[i] != "\n" {
			if task_str_slice[i] != "" {
				task_list[i] = ParseTask(task_str_slice[i])
				task_list[i].Index = i
			}
		}
	}

	return task_list
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
	tasks := TaskList()
	task := Task{
		Priority: 0,
		Content:  s,
		Date:     time.Now(),
		Done:     false,
		Index:    len(tasks) - 1,
	}
	return task
}

// Append a new task to the end of the tasks file
func WriteTask(task Task) error {
	f, err := os.OpenFile(json_path, os.O_APPEND|os.O_WRONLY, 0600)
	// check the error to see if we need to create a new file
	if err != nil {
		f2, nErr := os.Create(json_path)
		// if creation of new file fails, log it
		if nErr != nil {
			log.Fatal(nErr)
		}
		// assign the new file variable to the nil f
		f = f2
	}
	json, _ := json.Marshal(task)
	f.Write(json)
	f.WriteString("\n")
	return nil
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

// Print all tasks in tasks.txt
func PrintAllTasks() {
	task_list := TaskList()
	for i := range task_list {
		t := task_list[i]
		if t.Done == false {
			PrintTask(&t, 0)
		}
	}
}

// Marks a task as complete by setting the complete field to true in the JSON file
func CompleteTask(index int) {
	task_list := TaskList()
	task_list[index].Done = true
	os.Remove(json_path)
	for i := range task_list {
		WriteTask(task_list[i])
	}
	fmt.Printf("Task Marked as complete: %s\n", task_list[index].Content)
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
				WriteTask(task)
				fmt.Printf("Task is added: %s\n", task.Content)
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
