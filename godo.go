package main

import (
  "log"
  "encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type Task struct {
	Priority int64
	Content  string
	Date     time.Time
	Done     bool
	Index    int
}

// Takes a json string and converts it to a Task struct,(without an index)
func ParseTask(j string) Task {
  var task Task
  b := []byte(j)
  json.Unmarshal(b, &task)
  return task
}

// Returns an array of Tasks, with indices
func TaskList() []Task {
	file, _ := ioutil.ReadFile("tasks.json")
	file_str := string(file)
	task_str_slice := strings.Split(file_str, "\n")
  task_list := make([]Task, len(task_str_slice))
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
// task := Task {Priority: 0, Content: "get groceries, Date: time.Now(), Done: false}
// task.String //= "0       newtask 2014-03-14 22:22:47.875460951 -0600 MDT false 
func (t Task) String() string {
	return fmt.Sprintf("%d\t%s\t%s\t%t", t.Priority, t.Content, t.Date, t.Done)
}

// Print a Task
func (t *Task) Print() {
	println(t.String())
}

// Build a Task with Task.Content from string, with default values
func buildTask(s string) Task {
  tasks := TaskList()
	task := Task{
		Priority: 0,
		Content:  s,
		Date:     time.Now(),
		Done:     false,
    Index:    len(tasks)-1,
	}
	return task
}

// Append a new task to the end of the tasks file
func WriteTask(task Task) error {
	f, err := os.OpenFile("tasks.json", os.O_APPEND|os.O_WRONLY, 0600)
  // check the error to see if we need to create a new file
	if err != nil {
    f2, nErr := os.Create("tasks.json")
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

// Print all tasks in tasks.txt
func PrintTasks() {
	task_list := TaskList()
	for i := range task_list {
    t := task_list[i]
    if t.Done == false {
      t.Print()
    }
	}
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
				PrintTasks()
			},
		},
	}
	app.Run(os.Args)
}
