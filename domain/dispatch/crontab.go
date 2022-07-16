package dispatch

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/robfig/cron/v3"
	etcdv3 "go.etcd.io/etcd/client/v3"
	"tztask/utils/dlock"
)

type DCrontab struct {
	locker  dlock.Locker
	leader  bool
	waiting bool
	rwMux   sync.RWMutex
	cron    *cron.Cron
}

func newDCrontab(locker dlock.Locker) *DCrontab {
	c := &DCrontab{
		locker:  locker,
		leader:  false,
		waiting: false,
		rwMux:   sync.RWMutex{},
		cron:    cron.New(cron.WithSeconds()),
	}
	c.locker.SetBeforeUnlock(func() {
		c.rwMux.Lock()
		c.leader = false
		c.waiting = false
	})
	c.locker.SetAfterUnlock(func() {
		c.rwMux.Unlock()
	})
	_, err := c.cron.AddFunc("*/50 * * * * ?", c.cbLeader)
	if err != nil {
		log.Printf("error DCrontab cbLeader err: %v", err)
	}
	return c
}

func NewDCrontabByEtcd(etcd *etcdv3.Client, lockName string) *DCrontab {
	return newDCrontab(dlock.NewEtcd(etcd, fmt.Sprintf("crontab_%s", lockName)))
}

// Start 启动
func (c *DCrontab) Start() {
	c.cbLeader()
	c.cron.Start()
}

func (c *DCrontab) isLeader() bool {
	c.rwMux.RLock()
	defer c.rwMux.RUnlock()
	return c.leader
}

func (c *DCrontab) Close() {
	err := c.locker.UnLock(context.Background())
	if err != nil {
		log.Printf("error DCrontab close err: %v", err)
	}
}

func (c *DCrontab) cbLeader() {
	if c.waiting {
		return
	}
	c.waiting = true
	c.rwMux.Lock()
	defer c.rwMux.Unlock()
	err := c.locker.Lock(context.Background())
	if err != nil {
		log.Printf("DCrontab fail lock err: %v", err)
		c.leader = false
		c.waiting = false
		return
	}
	log.Println("c.leader ...")
	c.leader = true
}

func (c *DCrontab) AddFunc(spec string, fun func()) (cron.EntryID, error) {
	return c.cron.AddFunc(spec, func() {
		if c.isLeader() {
			fun()
		}
	})
}

func (c *DCrontab) RemoveFunc(id cron.EntryID) {
	c.cron.Remove(id)
}
