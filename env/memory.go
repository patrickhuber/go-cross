package env

import "fmt"

type memory struct {
	data map[string]string
}

func NewMemory() Environment {
	return NewMemoryWithMap(map[string]string{})
}

func NewMemoryWithMap(data map[string]string) Environment {
	return &memory{
		data: data,
	}
}

func (e *memory) Get(key string) string {
	value, ok := e.Lookup(key)
	if !ok {
		return ""
	}
	return value
}

func (e *memory) Lookup(key string) (string, bool) {
	value, ok := e.data[key]
	return value, ok
}

func (e *memory) Set(key, value string) error {
	e.data[key] = value
	return nil
}

func (e *memory) Delete(key string) error {
	delete(e.data, key)
	return nil
}

func (e *memory) Export() map[string]string {
	clone := make(map[string]string)
	for key, value := range e.data {
		clone[key] = value
	}
	return clone
}

func (e *memory) Environ() []string {
	list := []string{}
	for k, v := range e.data {
		list = append(list, fmt.Sprintf("%s=%s", k, v))
	}
	return list
}
