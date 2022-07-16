package repo_impl

import (
	"context"
	"database/sql"
	jsoniter "github.com/json-iterator/go"
	"log"
	"tztask/domain/entity"
	"tztask/domain/event"
	"tztask/utils/di"
)

const (
	taskTableSQL = `
CREATE TABLE IF NOT EXISTS "task"(
"name" VARCHAR(64) PRIMARY KEY  NOT NULL,
"expr" VARCHAR(64) NOT NULL,
"command" TEXT NOT NULL,
"status" SMALLINT NOT NULL,
"create_time" BIGINT NOT NULL,
"update_time" BIGINT NOT NULL)`

	taskGetSQL    = "SELECT expr,name,command,status,create_time,update_time FROM task WHERE name = ?"
	taskSaveSQL   = "INSERT INTO task (expr,name,command,status,create_time,update_time) VALUES (?,?,?,?,?,?) ON CONFLICT(name) DO UPDATE SET expr=?,command=?,status=?,update_time=?"
	taskLoadSQL   = "SELECT expr,name,command,status,create_time,update_time FROM task WHERE status!=-1"
	taskDeleteSQL = "UPDATE task SET status=-1 WHERE name = ?"
)

type TaskRepoSQLite struct {
	DB *sql.DB `di:"conf.SQLite"`
}

func NewTaskRepoSQLite() *TaskRepoSQLite {
	t := &TaskRepoSQLite{}
	di.MustBind(t)
	t.initTable()
	return t
}

func (t *TaskRepoSQLite) initTable() {
	_, err := t.DB.Exec(taskTableSQL)
	if err != nil {
		log.Printf("TaskRepoImpl.initTable err: %v", err)
	}
}

func (t *TaskRepoSQLite) Get(ctx context.Context, name string) (*entity.Task, error) {
	task := &entity.Task{}
	row := t.DB.QueryRowContext(ctx, taskGetSQL, name)
	var command string
	err := row.Scan(&task.Expr, &task.Name, &command, &task.Status, &task.CreateTime, &task.UpdateTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if command != "" {
		err := jsoniter.UnmarshalFromString(command, &task.Command)
		if err != nil {
			return nil, err
		}
	}
	return task, nil
}

func (t *TaskRepoSQLite) Load(ctx context.Context) ([]*entity.Task, error) {
	rows, err := t.DB.QueryContext(ctx, taskLoadSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tasks := make([]*entity.Task, 0)
	for rows.Next() {
		task := &entity.Task{}
		var command string
		err := rows.Scan(&task.Expr, &task.Name, &command, &task.Status, &task.CreateTime, &task.UpdateTime)
		if err != nil {
			return nil, err
		}
		if command != "" {
			err := jsoniter.UnmarshalFromString(command, &task.Command)
			if err != nil {
				return nil, err
			}
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (t *TaskRepoSQLite) Save(ctx context.Context, task *entity.Task) (*entity.Task, error) {
	cmd := task.Command.ToString()
	_, err := t.DB.ExecContext(ctx, taskSaveSQL,
		task.Expr,
		task.Name,
		cmd,
		task.Status,
		task.CreateTime,
		task.UpdateTime,
		task.Expr,
		cmd,
		task.Status,
		task.UpdateTime,
	)
	if err != nil {
		return nil, err
	}
	event.Trigger("dispatch.LoadTasks", task.Name)
	return task, nil
}

func (t *TaskRepoSQLite) Delete(ctx context.Context, name string) error {
	_, err := t.DB.ExecContext(ctx, taskDeleteSQL, name)
	if err != nil {
		return err
	}
	event.Trigger("dispatch.RemoveTask", name)
	return nil
}
