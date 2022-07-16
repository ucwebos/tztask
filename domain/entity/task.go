package entity

import (
	jsoniter "github.com/json-iterator/go"
)

type Task struct {
	Expr       string   `json:"expr" yaml:"expr"`
	Name       string   `json:"name" yaml:"name"`
	Command    *Command `json:"command" yaml:"command"`
	Status     int32    `json:"status"` // 1 正常 2 停止 -1 删除
	CreateTime int64    `json:"create_time"`
	UpdateTime int64    `json:"update_time"`
}

func (t *Task) ToString() string {
	str, _ := jsoniter.MarshalToString(t)
	return str
}

func (t *Task) FromJSON(buf []byte) error {
	return jsoniter.Unmarshal(buf, t)
}

type Command struct {
	Type   string `json:"type" yaml:"type"`
	Method string `json:"method" yaml:"method"`
	Target string `json:"target" yaml:"target"`
}

func (c *Command) ToString() string {
	str, _ := jsoniter.MarshalToString(c)
	return str
}

type TaskDashboard struct {
	TotalNum int     `json:"total_num"`
	TaskList []*Task `json:"task_list"`
}
