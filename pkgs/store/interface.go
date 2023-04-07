package store

const ErrNotFound = "No such key"

type Store interface {
	Set(key string, value string) error
	Get(key string) (string, error)
}
