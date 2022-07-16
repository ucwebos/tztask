package repo_impl

import (
	"context"
	"database/sql"
	"log"
	"tztask/domain/entity"
	"tztask/utils/di"
	"tztask/utils/sequence"
)

const (
	httpResultTableSQL = `
CREATE TABLE IF NOT EXISTS "http_result"(
"id" BIGINT PRIMARY KEY  NOT NULL,
"time" BIGINT NOT NULL,
"task_name" VARCHAR(64) NOT NULL,
"target" VARCHAR(255) NOT NULL,
"target_to" VARCHAR(255) NOT NULL,
"status_code" SMALLINT NOT NULL,
"content_length" BIGINT NOT NULL,
"content_type" VARCHAR(64) NOT NULL,
"raw" TEXT NOT NULL,
"title" VARCHAR(255) NOT NULL,
"description" TEXT NOT NULL)`

	httpResultSaveSQL = "INSERT INTO http_result (id,time,task_name,target,target_to,status_code,content_length,content_type,raw,title,description) VALUES (?,?,?,?,?,?,?,?,?,?,?)"

	httpResultQueryByTaskCountSQL = "SELECT count(1) FROM http_result WHERE task_name=?"

	httpResultQueryByTaskSQL = "SELECT id,time,task_name,target,target_to,status_code,content_length,content_type,raw,title,description FROM http_result WHERE task_name=? LIMIT ? OFFSET ?"
)

type HttpResultRepoSQLite struct {
	DB *sql.DB `di:"conf.SQLite"`
}

func NewHttpResultRepoSQLite() *HttpResultRepoSQLite {
	h := &HttpResultRepoSQLite{}
	di.MustBind(h)
	h.initTable()
	return h
}

func (h *HttpResultRepoSQLite) initTable() {
	_, err := h.DB.Exec(httpResultTableSQL)
	if err != nil {
		log.Printf("TaskRepoImpl.initTable err: %v", err)
	}
}

func (h *HttpResultRepoSQLite) Save(ctx context.Context, input *entity.HttpResult) error {
	id, err := sequence.ID()
	if err != nil {
		return err
	}
	_, err = h.DB.ExecContext(ctx, httpResultSaveSQL,
		id,
		input.Time,
		input.TaskName,
		input.Target,
		input.TargetTo,
		input.StatusCode,
		input.ContentLength,
		input.ContentType,
		input.Raw,
		input.Title,
		input.Description,
	)
	if err != nil {
		return err
	}
	return nil
}

func (h *HttpResultRepoSQLite) QueryByTask(ctx context.Context, taskName string, page int, size int) ([]*entity.HttpResult, int, error) {
	var total int
	row := h.DB.QueryRowContext(ctx, httpResultQueryByTaskCountSQL, taskName)
	if err := row.Scan(&total); err != nil {
		return nil, 0, err
	}
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * size
	rows, err := h.DB.QueryContext(ctx, httpResultQueryByTaskSQL, taskName, size, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	results := make([]*entity.HttpResult, 0)
	for rows.Next() {
		result := &entity.HttpResult{}
		err := rows.Scan(
			&result.ID,
			&result.Time,
			&result.TaskName,
			&result.Target,
			&result.TargetTo,
			&result.StatusCode,
			&result.ContentLength,
			&result.ContentType,
			&result.Raw,
			&result.Title,
			&result.Description,
		)
		if err != nil {
			return nil, 0, err
		}
		results = append(results, result)
	}
	return results, total, nil
}
