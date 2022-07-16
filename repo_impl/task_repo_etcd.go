package repo_impl

import (
	"context"
	"log"
	"strings"
	"tztask/conf"

	etcdv3 "go.etcd.io/etcd/client/v3"

	"tztask/domain/entity"
	"tztask/domain/repo"
	"tztask/utils/di"
)

const (
	taskKeyPrefix = "tztask/task/"
)

type TaskRepoEtcd struct {
	Etcd       *etcdv3.Client `di:"conf.Etcd"`
	SQLiteRepo repo.TaskRepo  `di:"repo_impl.TaskRepoSQLite"`
}

func NewTaskRepoEtcd() repo.TaskRepo {
	if !conf.App.Distributed {
		return NewTaskRepoSQLite()
	}
	t := &TaskRepoEtcd{}
	di.MustBind(t)
	go t.onNotify()
	return t
}

func (t *TaskRepoEtcd) key(name string) string {
	return taskKeyPrefix + name
}
func (t *TaskRepoEtcd) parseKey(key string) string {
	return strings.Replace(key, taskKeyPrefix, "", 1)
}

func (t *TaskRepoEtcd) Get(ctx context.Context, name string) (*entity.Task, error) {
	var (
		key  = t.key(name)
		task = &entity.Task{}
	)
	rs, err := t.Etcd.Get(ctx, key)
	if err != nil {
		log.Panicf("Etcd.Get err: %v", err)
		return nil, err
	}
	if len(rs.Kvs) > 0 {
		val := rs.Kvs[0].Value
		err = task.FromJSON(val)
		if err != nil {
			log.Panicf("Etcd.Get Unmarshal err: %v", err)
			return nil, err
		}
	}
	return task, nil
}

func (t *TaskRepoEtcd) Load(ctx context.Context) ([]*entity.Task, error) {
	tasks := make([]*entity.Task, 0)
	rs, err := t.Etcd.Get(ctx, taskKeyPrefix, etcdv3.WithPrefix())
	if err != nil {
		log.Panicf("Etcd.Get err: %v", err)
		return nil, err
	}
	for _, kv := range rs.Kvs {
		task := &entity.Task{}
		err = task.FromJSON(kv.Value)
		if err != nil {
			log.Panicf("Etcd.Get Unmarshal err: %v", err)
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (t *TaskRepoEtcd) Save(ctx context.Context, task *entity.Task) (*entity.Task, error) {
	var (
		key   = t.key(task.Name)
		value = task.ToString()
	)
	_, err := t.Etcd.Put(ctx, key, value)
	if err != nil {
		return nil, err
	}
	return task, err
}

func (t *TaskRepoEtcd) Delete(ctx context.Context, name string) error {
	_, err := t.Etcd.Delete(ctx, t.key(name))
	if err != nil {
		return err
	}
	return nil
}

func (t *TaskRepoEtcd) onNotify() {
	ch := t.Etcd.Watch(context.Background(), taskKeyPrefix, etcdv3.WithPrefix())
	for rs := range ch {
		for _, event := range rs.Events {
			switch event.Type {
			case etcdv3.EventTypeDelete:
				taskName := t.parseKey(string(event.Kv.Key))
				err := t.SQLiteRepo.Delete(context.Background(), taskName)
				if err != nil {
					log.Printf("SQLiteRepo.Delete %s,err:%v", taskName, err)
					continue
				}
			case etcdv3.EventTypePut:
				taskName := t.parseKey(string(event.Kv.Key))
				task := &entity.Task{}
				err := task.FromJSON(event.Kv.Value)
				if err != nil {
					log.Printf("EventTypePut %s,FromJSON err:%v", taskName, err)
					continue
				}
				_, err = t.SQLiteRepo.Save(context.Background(), task)
				if err != nil {
					log.Printf("SQLiteRepo.Save %s,err:%v", string(event.Kv.Value), err)
					continue
				}
			}
		}
	}
	t.onNotify()
}
