package di

import (
	"reflect"
)

type RegistryOption func(*registryItem)

type registryItem struct {
	value      interface{}
	hasBind    bool
	shouldBind bool
}

func newRegistryItem(value interface{}) *registryItem {
	result := &registryItem{
		value:      value,
		hasBind:    false,
		shouldBind: true,
	}

	t := reflect.TypeOf(value)
	for t.Kind() == reflect.Ptr || t.Kind() == reflect.Interface {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		result.shouldBind = false
	}
	return result
}

type registry struct {
	registry map[string]*registryItem
}

func NewRegistry() *registry {
	return &registry{
		registry: map[string]*registryItem{},
	}
}

func (r *registry) Register(name string, value interface{}, opts ...RegistryOption) {
	if r.get(name) != nil {
		panic("duplicate di register for name:" + name)
	}

	item := newRegistryItem(value)
	for _, o := range opts {
		o(item)
	}

	r.registry[name] = item
}

func (r *registry) FetchByName(name string) (interface{}, error) {
	resultItem := r.getByName(name)
	if resultItem == nil {
		return nil, nil
	}

	return r.bindItem(resultItem)
}

func (r *registry) FetchByType(t reflect.Type) (interface{}, error) {
	resultItem := r.getByType(t)

	if resultItem == nil {
		return nil, nil
	}

	return r.bindItem(resultItem)
}

func (r *registry) Bind(target interface{}) error {
	b := NewBinder(target, r)
	if err := b.Bind(); err != nil {
		return err
	}
	return nil
}

func (r *registry) bindItem(item *registryItem) (interface{}, error) {
	if !item.shouldBind {
		return item.value, nil
	}
	if item.hasBind {
		return item.value, nil
	}

	if err := r.Bind(item.value); err != nil {
		return nil, err
	}

	item.hasBind = true
	return item.value, nil
}

func (r *registry) getByName(name string) *registryItem {
	result := r.get(name)
	if result == nil {
		return nil
	}

	v := reflect.ValueOf(result.value)

	switch v.Kind() {
	case reflect.Func:
		return newRegistryItem(r.processFunc(v))
	default:
		return result
	}
}

// 获取所需对象, 不填充
func (r *registry) getByType(t reflect.Type) *registryItem {
	for _, item := range r.registry {
		v := reflect.ValueOf(item.value)

		switch v.Kind() {
		case reflect.Func:
			if r.FuncMatchType(v, t) {
				return newRegistryItem(r.processFunc(v))
			}
		default:
			if v.Type().AssignableTo(t) {
				return item
			}
		}
	}

	return nil
}

func (r *registry) FuncMatchType(v reflect.Value, target reflect.Type) bool {
	t := v.Type()
	if t.NumOut() != 1 {
		return false
	}

	return t.Out(0).AssignableTo(target)
}

func (r *registry) processFunc(v reflect.Value) interface{} {
	t := v.Type()
	input := make([]reflect.Value, 0)

	for i := 0; i != t.NumIn(); i++ {
		in := t.In(i)
		inputValue, err := r.FetchByType(in)
		if err != nil {
			return nil
		}

		input = append(input, reflect.ValueOf(inputValue))
	}

	result := v.Call(input)
	return result[0].Interface()
}

func (r *registry) get(name string) *registryItem {
	return r.registry[name]
}
