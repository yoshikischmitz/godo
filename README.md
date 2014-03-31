godo
======
A simple command line task tracking app to help me learn Go. Uses cli by codegangsta. By default it creates a tasks.json file in your home folder if you're on linux, or in My Documents on windows.

Usage:

    $ godo ls
    [1]     [2014-3-27]     Get groceries
    [2]     [2014-3-27]     Fix Issue #4501
      [1]   [2014-3-27]     Email client
      [2]   [2014-3-27]     Follow up with client
    [3]     [2014-3-28]     Add more features to Godo
    $ godo add "Update Godo readme file"
    Task is added: Update Godo readme file
    $ godo ls
    [1]     [2014-3-27]     Get groceries
    [2]     [2014-3-27]     Fix Issue #4501
      [1]   [2014-3-27]     Email client
      [2]   [2014-3-27]     Follow up with client
    [3]     [2014-3-28]     Add more features to Godo
    [4]     [2014-3-29]     Update Godo readme file
    $ godo subadd 2 "Close issue"
    $ godo complete 1
    Task Marked as complete: Get groceries
    $ godo ls
    [1]     [2014-3-27]     Fix Issue #4501
      [1]   [2014-3-27]     Email client
      [2]   [2014-3-27]     Follow up with client
      [3]   [2014-3-29]     Close Issue
    [2]     [2014-3-28]     Add more features to Godo
    [3]     [2014-3-29]     Update Godo readme file
    
Though I do use this myself, if you want a more stable and (currently)feature-rich solution I recommend the excellent bash based <a href-"todotxt.com">todo.txt</a>

Windows and Linux binaries are available <a href="https://sourceforge.net/projects/godo-cli/files/?source=navbar">here</a>.

Current features:
- Add new tasks with command godo add "todo text here"
- Add subtasks with godo subadd [parent task number] "subtask here"
- List tasks with godo ls
- Complete tasks with godo complete [task #]
- Keeps track of tasks in a simple JSON file

There are many planned features listed below. Along with everything mentioned in the Todo section below, my current focus is creating a "board" mode where you can navigate, edit, and add new tasks interactively while also staying in-sync with any changes you make from other shells. Also planned is the ability to expose the todos as a simple web service for group task-tracking. The idea is for the interface to live as both a git-style cli and a vim-style editor.

Any and all ideas and suggestions are welcome!

TODO:
- Add config options for tasks.json location
- Assign priorities to task
- Display tasks by priority
- Tags, Locations, Filters to tasks
- Add flag to ls to display completed items as well
- Add muiltask completion functionality
- Add shortcuts for all commands
