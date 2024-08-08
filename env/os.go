package env

import (
	"os"
	"strings"
)

type env struct {
}

func New() Environment {
	return &env{}
}

func (e *env) Get(key string) string {
	return os.Getenv(key)
}

func (e *env) Set(key, value string) error {
	return os.Setenv(key, value)
}

func (e *env) Lookup(key string) (string, bool) {
	return os.LookupEnv(key)
}

func (e *env) Delete(key string) error {
	return os.Unsetenv(key)
}

func (e *env) Export() map[string]string {
	clone := make(map[string]string)
	for _, item := range os.Environ() {

		split := strings.Split(item, "=")
		if len(split) != 2 {
			continue
		}

		key := split[0]
		value := split[1]

		clone[key] = value
	}
	return clone
}

func (e *env) Environ() []string {
	return os.Environ()
}
