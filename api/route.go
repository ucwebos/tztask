package api

import (
	"net/http"
	"tztask/conf"
	"tztask/utils/simple_server"
)

func routes() *simple_server.Router {
	r := simple_server.NewRouter()
	h := NewHandler()

	r.POST("/task/add", h.addTask)
	r.POST("/task/set", h.setTask)
	r.POST("/task/delete", h.deleteTask)
	// 任务结果查询接口
	r.GET("/task/result", h.taskResult)
	// 查看系统现有任务
	r.GET("/", h.dashboard)

	return r
}

func HTTPServe() error {
	s := &simple_server.Server{}
	s.SetRoute(routes())
	return http.ListenAndServe(conf.App.HttpAddr, s)
}
