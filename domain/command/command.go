package command

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
	"strings"
	"tztask/domain/entity"
	"tztask/utils/di"
)

const (
	TypeHttp = "http"
)

type Command interface {
	ID() string
	Spec() string
	Verify() error
	TryNum() int
	Run() error
}

func Parse(task *entity.Task) (Command, error) {
	if task.Command == nil {
		return nil, errors.New("task illegal")
	}
	cmd, err := parseTo(task)
	if err != nil {
		return nil, errors.Wrapf(err, "task parse err:")
	}
	if err := cmd.Verify(); err != nil {
		return nil, errors.Wrapf(err, "Command Verify:")
	}

	return cmd, nil
}

func parseTo(task *entity.Task) (Command, error) {
	var (
		tmp  = strings.Split(task.Expr, " ")
		expr = task.Expr
	)

	if len(tmp) == 5 {
		expr = fmt.Sprintf("0 %s", expr)
	}

	switch task.Command.Type {
	case TypeHttp:
		cmd := &HttpCommand{
			Expr:    expr,
			Name:    task.Name,
			Command: task.Command,
		}
		di.MustBind(cmd)
		return cmd, nil
	default:
		return nil, errors.Errorf("task.Command.Type [%s] not allowd ", task.Command.Type)
	}
}

func RetryRun(cmd Command) {
	tryNum := cmd.TryNum()
	if tryNum == 0 {
		tryNum = 1
	}
	log.Printf("Run %s \n", cmd.ID())
	for i := 0; i < tryNum; i++ {
		err := cmd.Run()
		if err == nil {
			break
		}
		log.Printf("cmd.Run err: %v", err)
	}
}
