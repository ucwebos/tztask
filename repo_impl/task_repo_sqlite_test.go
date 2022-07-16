package repo_impl

import (
	"context"
	"fmt"
	"testing"
	"time"

	jsoniter "github.com/json-iterator/go"

	"tztask/domain/entity"
)

func TestTaskRepoSQLite_Save(t *testing.T) {
	_, err := NewTaskRepoSQLite().Save(context.Background(), &entity.Task{
		Expr: "* * * * * *",
		Name: "job-1",
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
}

func TestTaskRepoSQLite_Get(t *testing.T) {
	task, err := NewTaskRepoSQLite().Get(context.Background(), "job-2")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(jsoniter.MarshalToString(task))
}

func TestTaskRepoSQLite_Load(t *testing.T) {
	tasks, err := NewTaskRepoSQLite().Load(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(jsoniter.MarshalToString(tasks))
}

func TestTaskRepoSQLite_Delete(t *testing.T) {
	err := NewTaskRepoSQLite().Delete(context.Background(), "job-1")
	if err != nil {
		t.Fatal(err)
	}
}
