package storage

type Storage map[string]string

// NewStorage ...
func NewStorage() Storage {
	return make(Storage)
}
