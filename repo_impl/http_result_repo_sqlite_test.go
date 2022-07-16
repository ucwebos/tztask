package repo_impl

import (
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"testing"
	"time"
	"tztask/domain/entity"
	"tztask/utils/sequence"
)

func TestHttpResultRepoSQLite_Save(t *testing.T) {
	id, _ := sequence.ID()
	err := NewHttpResultRepoSQLite().Save(context.Background(), &entity.HttpResult{
		ID:            id,
		Time:          time.Now().Unix(),
		TaskName:      "job-2",
		Target:        "http://baidu.com",
		TargetTo:      "http://www.baidu.com",
		StatusCode:    200,
		ContentLength: 210200,
		ContentType:   "text/html",
		Raw:           "xxxxx",
		Title:         "xx-title",
		Description:   "xx-description",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestHttpResultRepoSQLite_QueryByTask(t *testing.T) {
	results, total, err := NewHttpResultRepoSQLite().QueryByTask(context.Background(), "job-1", 1, 100)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(total)
	fmt.Println(jsoniter.MarshalToString(results))
}
