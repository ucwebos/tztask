package event

var events = map[string]func(arg string){}

func RegisterEvent(name string, f func(arg string)) {
	events[name] = f
}

func Trigger(name string, arg string) {
	if f, ok := events[name]; ok {
		f(arg)
	}
}
