package repo_impl

import (
	"testing"
	"tztask/conf"
	"tztask/utils/sequence"
)

func TestMain(m *testing.M) {
	conf.App = &conf.AppConfig{
		ServerName: "tztask",
		HttpAddr:   "8067",
		SQLite:     &conf.SQLite{DBFile: "test.db"},
		//Etcd: &conf.Etcd{
		//	Endpoints: []string{"127.0.0.1:2379"},
		//	Timeout:   5,
		//},
		ShardNum: 5,
	}
	conf.DIRegister()
	sequence.Init()
	DIRegister()
	m.Run()
}
