package http_io

import "tztask/domain/entity"

type AddTaskReq struct {
	Expr    string          `json:"expr" yaml:"expr"`
	Name    string          `json:"name" yaml:"name"`
	Command *entity.Command `json:"command" yaml:"command"`
}

type AddTaskResp struct {
	Success bool `json:"success"`
}

type SetTaskReq struct {
	Expr    string          `json:"expr" yaml:"expr"`
	Name    string          `json:"name" yaml:"name"`
	Command *entity.Command `json:"command" yaml:"command"`
}

type SetTaskResp struct {
	Success bool `json:"success"`
}

type DeleteTaskReq struct {
	Name string `json:"name" yaml:"name"`
}

type DeleteTaskResp struct {
	Success bool `json:"success"`
}

type TaskResultResp struct {
	List  []*entity.HttpResult `json:"list"`
	Total int                  `json:"total"`
}
