package store

type InMemStore interface {
	Get(key string) (string, bool)
	Set(key, value string)
	Delete(key string)
}
