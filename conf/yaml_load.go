package conf

import (
	"io/ioutil"
	"log"
	"reflect"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"tztask/utils"
)

func YamlLoad(yamlFile string, rs interface{}) error {
	if yamlFile == "" {
		return errors.New("yamlFile is empty")
	}
	if rs == nil {
		return errors.New("param rs is nil")
	}
	typ := reflect.TypeOf(rs)
	if typ.Kind() != reflect.Ptr {
		return errors.New("cannot apply to non-pointer struct")
	}
	if utils.Exists(yamlFile) {
		buf, _ := ioutil.ReadFile(yamlFile)
		err := yaml.Unmarshal(buf, rs)
		if err != nil {
			log.Printf("%s yaml file parse err:%v", yamlFile, err)
			return err
		}
	}
	return nil
}
