package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"time"
	"github.com/nsf/termbox-go"
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
		os.Open(json_path)
	}
	tasks_root := TasksRoot{}
	json.Unmarshal(file, &tasks_root)
	return tasks_root
}

// Build a Task with Task.Content from string with default values.
func buildTask(s string) Task {
	task := Task{
		Priority: 0,
		Content:  s,
		Date:     time.Now(),
		Done:     false,
		Index:    len(task_root.Tasks),
	}
	return task
}

func WriteJson() {
	json, _ := json.MarshalIndent(task_root, "", "  ")
	os.Remove(json_path)
	ioutil.WriteFile(json_path, json, 0600)
}

// Add a new task to tasks file
func AddTask(task Task) error {
	task_root.Tasks = append(task_root.Tasks, task)
	WriteJson()
	return nil
}

func AddSubTask(task *Task, index int) {
	i, _ := real_index(index)
	t := &task_root.Tasks[i]
	t.SubTasks = append(t.SubTasks, task)
	WriteJson()
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

func real_index(pidx int) (int, error) {
	// Return print index to 0-index
	pidx -= 1
	// i keeps track of the real index in data,
	// j keeps track of the index of completed items, as seen by user
	var j int
	for i := range task_root.Tasks {
		if task_root.Tasks[i].Done == false {
			if j == pidx {
				return i, nil
			}
			j += 1
		}
	}
	return 0, errors.New("can't finx index")
}

// Sets the complete field to true at relative index for uncomplete tasks
func CompleteTask(index int) {
	i, _ := real_index(index)
	task_root.Tasks[i].Done = true
	fmt.Printf("Task Marked as complete: %s\n", task_root.Tasks[i].Content)
	WriteJson()

}

func ParseIndex(c *cli.Context) int {
	i, _ := strconv.ParseInt(c.Args().First(), 10, 0)
	return int(i)
}

func draw_string(s string, row int) {
	const coldef = termbox.ColorDefault
	for i, x := range s {
		termbox.SetCell(i, row, x, coldef, coldef)
	}
}

func DrawTask(t *Task, index int, row int) int {
	year, month, day := t.Date.Date()
	draw_string(
		fmt.Sprintf("[%d]\t[%d-%d-%d]\t%s\n", index, year, month, day, t.Content),
		row)
	rowdiff := 1
	for i, task := range t.SubTasks {
		i++
		rowdiff += DrawTask(task, i, row+rowdiff)
	}
	return rowdiff
}

// Print all tasks in tasks.json
func DrawAllTasks() {
	var j int
	row := 1
	for i := range task_root.Tasks {
		t := task_root.Tasks[i]
		if t.Done == false {
			j += 1
			row += DrawTask(&t, j, row)
		}
	}
}

func draw_tasks() {
	const coldef = termbox.ColorDefault
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	PrintAllTasks()
	termbox.Flush()
}



func runTermbox() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	draw_tasks()
loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break loop
			}
		case termbox.EventError:
			panic(ev.Err)
		}
		draw_tasks()
	}
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
				index := ParseIndex(c)
				CompleteTask(index)
			},
		},
		{
			Name: "board",
			Usage: "starts godo in go board",
			Action: func(c *cli.Context){
				runTermbox()
			},
		},
	}
	app.Run(os.Args)
}
