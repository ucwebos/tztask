package dispatch

import (
	"fmt"
	"testing"
	"tztask/conf"
	"tztask/domain/entity"
)

func TestNewManager(t *testing.T) {
	m := NewManager("", nil, 3)
	m.Start()

	str := "cdewfwefewfewfew"
	fmt.Println(len(str))
	fmt.Println(len([]byte(str)))

	task := &entity.Task{
		Expr: "*/2 * * * * *",
		Name: "test1",
		Command: &entity.Command{
			Type:   "http",
			Method: "GET",
			Target: "http://baidu.com/",
		},
	}
	m.AddTask(task)

	select {}
}

func TestNewManagerDistributed(t *testing.T) {
	etcdCfg := conf.Etcd{
		Endpoints: []string{"127.0.0.1:2379"},
		Timeout:   5, // s
	}
	etcd, err := etcdCfg.CreateInstance()
	if err != nil {
		t.Fatal(err)
	}
	m := NewManager("tztask", etcd, 3)
	m.Start()

	select {}
}
