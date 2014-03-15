package main 

import "testing"

func TestParse(t *testing.T) {
  json := `{"Priority":0,"Content":"hello world","Date":"2014-03-15T13:57:51.187311327-06:00","Done":false,"Index":0}`
  task := ParseTask(json)
  if task.Priority != 0 {
    t.Errorf("priority should be 0")
  }
  if task.Content != "hello world" {
    t.Errorf("content should be hello world")
  }
}

