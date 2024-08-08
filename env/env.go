package env

type Environment interface {
	Get(key string) string
	Set(key string, value string) error
	Lookup(key string) (string, bool)
	Export() map[string]string
	Environ() []string
	Delete(key string) error
}
