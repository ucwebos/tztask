package service

import (
	"context"
	"tztask/domain/entity"
	"tztask/domain/repo"
)

type HttpResultService struct {
	HttpResultRepo repo.HttpResultRepo `di:"repo_impl.HttpResultRepoSQLite"`
}

func NewHttpResultService() *HttpResultService {
	return &HttpResultService{}
}

func (s *HttpResultService) QueryByTask(ctx context.Context, taskName string, page, size int) ([]*entity.HttpResult, int, error) {
	return s.HttpResultRepo.QueryByTask(ctx, taskName, page, size)
}
