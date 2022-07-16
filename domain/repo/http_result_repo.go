package repo

import (
	"context"

	"tztask/domain/entity"
)

type HttpResultRepo interface {
	Save(ctx context.Context, input *entity.HttpResult) error
	QueryByTask(ctx context.Context, taskName string, page int, size int) ([]*entity.HttpResult, int, error)
}
