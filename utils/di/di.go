package di

import (
	"fmt"
	"reflect"
	"sync"
)

var (
	once      sync.Once
	gInstance *diInstance
)

type diInstance struct {
	registry *registry
}

func newInstance() *diInstance {
	return &diInstance{
		registry: NewRegistry(),
	}
}

func init() {
	once.Do(func() {
		gInstance = newInstance()
	})
}

func Register(name string, value interface{}) {
	gInstance.registry.Register(name, value)
}

func Get(name string) interface{} {
	return gInstance.registry.get(name)
}

func GetInst(name string) interface{} {
	item := gInstance.registry.get(name)
	if item == nil {
		return nil
	}
	return item.value
}

func Bind(target interface{}) error {
	return gInstance.registry.Bind(target)
}

func FetchByName(name string, target interface{}) error {
	if target == nil {
		return fmt.Errorf("can not bind to nil interface")
	}

	result, err := gInstance.registry.FetchByName(name)
	if err != nil {
		return err
	}
	return setResult(target, result)
}

func Fetch(target interface{}) error {
	if target == nil {
		return fmt.Errorf("can not bind to nil interface")
	}

	if reflect.TypeOf(target).Kind() != reflect.Ptr {
		return fmt.Errorf("should only bind to pointer")
	}

	result, err := gInstance.registry.FetchByType(reflect.TypeOf(target).Elem())
	if err != nil {
		return err
	}
	return setResult(target, result)
}

func MustFetchByName(name string, target interface{}) {
	if err := FetchByName(name, target); err != nil {
		panic(err)
	}
}

func MustFetch(target interface{}) {
	if err := Fetch(target); err != nil {
		panic(err)
	}
}

func MustBind(target interface{}) {
	if err := Bind(target); err != nil {
		panic(err)
	}
}

func MustBindALL() {
	for _, v := range gInstance.registry.registry {
		MustBind(v)
	}
}

func setResult(target interface{}, result interface{}) error {
	resultValue := reflect.ValueOf(result)
	v := reflect.ValueOf(target)
	for {
		if resultValue.Type().AssignableTo(v.Type()) && v.CanSet() {
			v.Set(resultValue)
			return nil
		}

		if (v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface) && !v.IsNil() {
			v = v.Elem()
		} else if (resultValue.Kind() == reflect.Ptr || resultValue.Kind() == reflect.Interface) && !resultValue.IsNil() {
			resultValue = resultValue.Elem()
		} else {
			return fmt.Errorf("can not assign %s to %s", resultValue.Type(), reflect.TypeOf(target))
		}
	}
}
