package conf

import "tztask/domain/entity"

var Jobs = make([]*entity.Task, 0)

type JobsConfig struct {
	Jobs []*entity.Task `yaml:"jobs"`
}

func LoadJobs() error {
	jcs := &JobsConfig{}
	err := YamlLoad(App.JobsFile, &jcs)
	if err != nil {
		return err
	}
	Jobs = jcs.Jobs
	return nil
}
