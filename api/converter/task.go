package converter

import (
	"time"
	"tztask/api/http_io"
	"tztask/domain/entity"
)

func FromAddTaskReq(req *http_io.AddTaskReq) *entity.Task {
	return &entity.Task{
		Expr:       req.Expr,
		Name:       req.Name,
		Command:    req.Command,
		Status:     1,
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now().Unix(),
	}
}

func FromSetTaskReq(task *entity.Task, req *http_io.SetTaskReq) (*entity.Task, bool) {
	change := false
	if req.Expr != "" && task.Expr != req.Expr {
		task.Expr = req.Expr
		change = true
	}
	if req.Command != nil && task.Command.ToString() != req.Command.ToString() {
		task.Command = req.Command
		change = true
	}
	return task, change
}
