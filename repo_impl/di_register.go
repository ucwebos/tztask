package repo_impl

import (
	"tztask/utils/di"
)

func DIRegister() {
	di.Register("repo_impl.TaskRepoSQLite", NewTaskRepoSQLite())
	di.Register("repo_impl.TaskRepoEtcd", NewTaskRepoEtcd())
	di.Register("repo_impl.HttpResultRepoSQLite", NewHttpResultRepoSQLite())
}
