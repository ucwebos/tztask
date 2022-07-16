package service

import (
	"tztask/utils/di"
)

func DIRegister() {
	di.Register("service.TaskService", NewTaskService())
	di.Register("service.HttpResultService", NewHttpResultService())
}
