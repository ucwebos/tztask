package api

import (
	"strconv"
	"tztask/api/converter"
	"tztask/api/http_io"
	"tztask/domain/command"
	"tztask/domain/service"
	"tztask/utils/di"
	"tztask/utils/simple_server"
)

type Handler struct {
	TaskService       *service.TaskService       `di:"service.TaskService"`
	HttpResultService *service.HttpResultService `di:"service.HttpResultService"`
}

func NewHandler() *Handler {
	h := &Handler{}
	di.MustBind(h)
	return h
}

func (h *Handler) addTask(ctx *simple_server.Context) {
	var (
		req  = &http_io.AddTaskReq{}
		resp = &http_io.AddTaskResp{}
	)
	if err := http_io.BindBody(ctx, &req); err != nil {
		http_io.JSONError(ctx, http_io.ErrParams)
		return
	}
	task := converter.FromAddTaskReq(req)
	if _, err := command.Parse(task); err != nil {
		http_io.JSONError(ctx, http_io.ErrParams)
		return
	}
	if err := h.TaskService.AddTask(ctx, task); err != nil {
		http_io.JSONError(ctx, http_io.ErrSystem)
		return
	}

	resp.Success = true
	http_io.JSONSuccess(ctx, resp)
}

func (h *Handler) setTask(ctx *simple_server.Context) {
	var (
		req  = &http_io.SetTaskReq{}
		resp = &http_io.SetTaskResp{}
	)
	if err := http_io.BindBody(ctx, &req); err != nil {
		http_io.JSONError(ctx, http_io.ErrParams)
		return
	}
	if req.Name == "" {
		http_io.JSONError(ctx, http_io.ErrParams)
		return
	}

	task, err := h.TaskService.GetTask(ctx, req.Name)
	if err != nil {
		http_io.JSONError(ctx, http_io.ErrSystem)
		return
	}
	if task == nil {
		_err := http_io.ErrUser
		_err.Err = "task not found"
		http_io.JSONError(ctx, _err)
		return
	}
	task, change := converter.FromSetTaskReq(task, req)
	if change {
		if err := h.TaskService.SetTask(ctx, task); err != nil {
			http_io.JSONError(ctx, http_io.ErrSystem)
			return
		}
	}
	resp.Success = true
	http_io.JSONSuccess(ctx, resp)
}

func (h *Handler) deleteTask(ctx *simple_server.Context) {
	var (
		req  = &http_io.DeleteTaskReq{}
		resp = &http_io.DeleteTaskResp{}
	)
	if err := http_io.BindBody(ctx, &req); err != nil {
		http_io.JSONError(ctx, http_io.ErrParams)
		return
	}
	if req.Name == "" {
		http_io.JSONError(ctx, http_io.ErrParams)
		return
	}
	if err := h.TaskService.DeleteTask(ctx, req.Name); err != nil {
		http_io.JSONError(ctx, http_io.ErrSystem)
		return
	}

	resp.Success = true
	http_io.JSONSuccess(ctx, resp)
}

func (h *Handler) dashboard(ctx *simple_server.Context) {
	stats, err := h.TaskService.TaskDashboard(ctx)
	if err != nil {
		http_io.JSONError(ctx, http_io.ErrSystem)
		return
	}
	http_io.JSONSuccess(ctx, stats)

}

func (h *Handler) taskResult(ctx *simple_server.Context) {
	var (
		taskName   = ctx.Param("name")
		_page      = ctx.Param("page")
		_size      = ctx.Param("size")
		page, size = 1, 100
		err        error
		resp       = &http_io.TaskResultResp{}
	)
	if taskName == "" {
		http_io.JSONError(ctx, http_io.ErrParams)
		return
	}
	if _page != "" {
		page, err = strconv.Atoi(_page)
		if err != nil {
			http_io.JSONError(ctx, http_io.ErrParams)
			return
		}
	}

	if _size != "" {
		size, err = strconv.Atoi(_size)
		if err != nil {
			http_io.JSONError(ctx, http_io.ErrParams)
			return
		}
	}

	list, total, err := h.HttpResultService.QueryByTask(ctx, taskName, page, size)
	if err != nil {
		http_io.JSONError(ctx, http_io.ErrParams)
		return
	}
	resp.List = list
	resp.Total = total
	http_io.JSONSuccess(ctx, resp)
}
