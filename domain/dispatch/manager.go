package dispatch

import (
	"context"
	"fmt"
	"log"
	"sync"
	"tztask/domain/entity"
	"tztask/domain/event"
	"tztask/domain/repo"
	"tztask/utils/di"

	"github.com/robfig/cron/v3"
	etcdv3 "go.etcd.io/etcd/client/v3"

	"tztask/conf"
	"tztask/domain/command"
	"tztask/utils"
)

var mg *Manager

type cronItem struct {
	cid  cron.EntryID
	expr string
	cmd  *entity.Command
}

type Manager struct {
	Distributed bool
	ShardNum    int
	cron        *cron.Cron
	DShards     map[uint32]*DCrontab
	cronMap     map[string]cronItem
	TaskRepo    repo.TaskRepo `di:"repo_impl.TaskRepoSQLite"`
	Mux         sync.Mutex
}

func NewManager(serverName string, etcd *etcdv3.Client, shardNum int) *Manager {
	if etcd == nil {
		m := &Manager{
			Distributed: false,
			cron:        cron.New(cron.WithSeconds()),
			cronMap:     map[string]cronItem{},
			Mux:         sync.Mutex{},
		}
		di.MustBind(m)
		return m
	}
	if shardNum == 0 {
		shardNum = 5
	}
	m := &Manager{
		Distributed: true,
		ShardNum:    shardNum,
		DShards:     make(map[uint32]*DCrontab, shardNum),
		cronMap:     map[string]cronItem{},
		Mux:         sync.Mutex{},
	}
	for i := 0; i < shardNum; i++ {
		dShard := NewDCrontabByEtcd(etcd, fmt.Sprintf("%s_%d", serverName, i))
		m.DShards[uint32(i)] = dShard
	}
	di.MustBind(m)
	m.LoadTasks("init")
	return m
}

func (m *Manager) LoadTasks(arg string) {
	tasks, err := m.TaskRepo.Load(context.Background())
	if err != nil {
		log.Printf("dispatch TaskRepo.Load err: %v", err)
		return
	}
	for _, task := range tasks {
		// 是否已存在
		if item, ok := m.cronMap[task.Name]; ok {
			if item.cmd.ToString() != task.Command.ToString() || task.Expr != item.expr {
				m.RemoveTask(task.Name)
				if err != nil {
					log.Printf("dispatch LoadTasks AddTask err: %v", err)
				}
			}
		}
		err = m.AddTask(task)
		if err != nil {
			log.Printf("dispatch LoadTasks AddTask err: %v", err)
		}
	}
}

func (m *Manager) AddTask(task *entity.Task) error {
	m.Mux.Lock()
	defer m.Mux.Unlock()
	cmd, err := command.Parse(task)
	if err != nil {
		return err
	}
	// 是否是分布式
	if m.Distributed {
		shard := utils.MurmurHash2(cmd.ID()) % uint32(m.ShardNum)
		dShard := m.DShards[shard]
		cronID, err := dShard.AddFunc(cmd.Spec(), func() {
			fmt.Println("xxxxxxx")
			command.RetryRun(cmd)
		})
		if err != nil {
			return err
		}
		m.cronMap[task.Name] = cronItem{
			cid:  cronID,
			expr: task.Expr,
			cmd:  task.Command,
		}
		return err
	}
	cronID, err := m.cron.AddFunc(cmd.Spec(), func() {
		command.RetryRun(cmd)
	})
	m.cronMap[task.Name] = cronItem{
		cid:  cronID,
		expr: task.Expr,
		cmd:  task.Command,
	}
	return err
}

func (m *Manager) RemoveTask(taskName string) {
	m.Mux.Lock()
	defer m.Mux.Unlock()
	if item, ok := m.cronMap[taskName]; ok {
		m.cron.Remove(item.cid)
	}
}

func (m *Manager) Start() {
	if !m.Distributed {
		go m.cron.Start()
		return
	}
	for _, crontab := range m.DShards {
		go crontab.Start()
	}
}

func (m *Manager) Close() {
	if !m.Distributed {
		m.cron.Stop()
	}
	for _, crontab := range m.DShards {
		crontab.Close()
	}
}

func Start() error {
	if conf.App.Distributed {
		mg = NewManager(conf.App.ServerName, conf.GetEtcd(), conf.App.ShardNum)
	} else {
		mg = NewManager(conf.App.ServerName, nil, 0)
	}
	// 注册load事件
	event.RegisterEvent("dispatch.LoadTasks", mg.LoadTasks)
	// 注册remove事件
	event.RegisterEvent("dispatch.RemoveTask", mg.RemoveTask)
	// 启动
	mg.Start()
	return nil
}

func Close() {
	mg.Close()
}
