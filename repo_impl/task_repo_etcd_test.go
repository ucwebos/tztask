package repo_impl

import (
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"testing"
	"time"
	"tztask/domain/entity"
)

func TestTaskRepoEtcd_Save(t *testing.T) {
	_, err := NewTaskRepoEtcd().Save(context.Background(), &entity.Task{
		Expr: "* * * * * *",
		Name: "job-4",
		Command: &entity.Command{
			Type:   "http",
			Method: "GET",
			Target: "http://baidu.com",
		},
		Status:     1,
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now().Unix(),
	})
	if err != nil {
		t.Fatal(err)
	}

	select {
	case <-time.After(3 * time.Second):

	}
}

func TestTaskRepoEtcd_Get(t *testing.T) {
	task, err := NewTaskRepoEtcd().Get(context.Background(), "job-4")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(jsoniter.MarshalToString(task))
	select {
	case <-time.After(3 * time.Second):

	}
}

func TestTaskRepoEtcd_Load(t *testing.T) {
	tasks, err := NewTaskRepoEtcd().Load(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(jsoniter.MarshalToString(tasks))
	select {
	case <-time.After(3 * time.Second):

	}
}

func TestTaskRepoEtcd_Delete(t *testing.T) {
	err := NewTaskRepoEtcd().Delete(context.Background(), "job-4")
	if err != nil {
		t.Fatal(err)
	}
	select {
	case <-time.After(3 * time.Second):

	}
}
