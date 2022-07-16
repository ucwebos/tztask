package repo

import (
	"context"

	"tztask/domain/entity"
)

type TaskRepo interface {
	Get(ctx context.Context, name string) (*entity.Task, error)
	Load(ctx context.Context) ([]*entity.Task, error)
	Save(ctx context.Context, input *entity.Task) (*entity.Task, error)
	Delete(ctx context.Context, name string) error
}
