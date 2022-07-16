package conf

import (
	etcdv3 "go.etcd.io/etcd/client/v3"
	"log"
	"os"
	"path/filepath"
	"tztask/utils"
	"tztask/utils/di"
)

const configFile = "config.yaml"

var (
	App = &AppConfig{}
)

type AppConfig struct {
	ServerName  string  `json:"server_name" yaml:"server_name"`
	HttpAddr    string  `json:"http_addr" yaml:"http_addr"`
	SQLite      *SQLite `json:"SQLite" yaml:"SQLite"`
	Etcd        *Etcd   `json:"etcd" yaml:"etcd"`
	JobsFile    string  `json:"jobs_file" yaml:"jobs_file"`
	ShardNum    int     `json:"shard_num" yaml:"shard_num"`
	Distributed bool
}

func Init() {
	err := YamlLoad(configFile, App)
	if err != nil {
		panic(err)
	}
	if App.JobsFile != "" {
		absPath, err := filepath.Abs(App.JobsFile)
		if err != nil {
			log.Panicf("jobs_file err: %v", err)
		}
		if !utils.Exists(absPath) {
			os.Create(absPath)
		}
		App.JobsFile = absPath
	}
	if App.Etcd != nil {
		App.Distributed = true
	}
	LoadJobs()
	DIRegister()
}

func DIRegister() {
	if App.Etcd != nil {
		instEtcd, err := App.Etcd.CreateInstance()
		if err != nil {
			log.Panicf("init DI conf.Etcd] err: %v", err)
		}
		di.Register("conf.Etcd", instEtcd)
	}
	instSQLite, err := App.SQLite.CreateInstance()
	if err != nil {
		log.Panicf("init DI conf.SQLLite] err: %v", err)
	}
	di.Register("conf.SQLite", instSQLite)
}

func GetEtcd() *etcdv3.Client {
	inst := di.GetInst("conf.Etcd")
	if inst == nil {
		return nil
	}
	return inst.(*etcdv3.Client)
}
