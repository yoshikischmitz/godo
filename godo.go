package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"os"
	"strconv"
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

// Takes a string and converts it to a Task struct,(without an index)
func ParseTask(s string) Task {
	task_str := strings.Split(s, "\t")
	task_priority, _ := strconv.ParseInt(task_str[0], 2, 0)
	task_time, _ := time.Parse(task_str[2], task_str[2])
	task_done, _ := strconv.ParseBool(task_str[3])
	task := Task{
		Priority: task_priority,
		Content:  task_str[1],
		Date:     task_time,
		Done:     task_done,
	}
	return task
}

// Returns an array of Tasks, with indices
func TaskList() []Task {
	file, _ := ioutil.ReadFile("tasks.txt")
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
	task := Task{
		Priority: 0,
		Content:  s,
		Date:     time.Now(),
		Done:     false,
	}
	return task
}

// Append a new task to the end of the tasks file
func WriteTask(task Task) error {
	text := task.String()
	f, err := os.OpenFile("tasks.txt", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	f.WriteString("\n" + text)
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
