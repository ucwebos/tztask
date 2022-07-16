package service

import (
	"context"
	"log"
	"tztask/conf"
	"tztask/domain/command"
	"tztask/utils/di"

	"tztask/domain/entity"
	"tztask/domain/repo"
)

type TaskService struct {
	TaskRepo repo.TaskRepo `di:"repo_impl.TaskRepoEtcd"`
}

func NewTaskService() *TaskService {
	s := &TaskService{}
	di.MustBind(s)

	// 文件配置入DB
	s.loadFileJobs()
	return s
}

func (s *TaskService) loadFileJobs() {
	for _, task := range conf.Jobs {
		_, err := command.Parse(task)
		if err != nil {
			log.Printf("conf jobs[%s] parse err: %v", task.Name, err)
			continue
		}
		task.Status = 1
		err = s.AddTask(context.Background(), task)
		if err != nil {
			log.Printf("conf jobs AddTask err: %v", err)
		}
	}
}

func (s *TaskService) AddTask(ctx context.Context, task *entity.Task) error {
	_, err := s.TaskRepo.Save(ctx, task)
	if err != nil {
		return err
	}
	return nil
}

func (s *TaskService) GetTask(ctx context.Context, taskName string) (*entity.Task, error) {
	return s.TaskRepo.Get(ctx, taskName)
}

func (s *TaskService) SetTask(ctx context.Context, task *entity.Task) error {
	_, err := s.TaskRepo.Save(ctx, task)
	return err
}

func (s *TaskService) TaskDashboard(ctx context.Context) (*entity.TaskDashboard, error) {
	tasks, err := s.TaskRepo.Load(ctx)
	if err != nil {
		return nil, err
	}
	return &entity.TaskDashboard{
		TotalNum: len(tasks),
		TaskList: tasks,
	}, nil
}

func (s *TaskService) DeleteTask(ctx context.Context, taskName string) error {
	return s.TaskRepo.Delete(ctx, taskName)
}
