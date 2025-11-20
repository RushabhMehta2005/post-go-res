package wal

import "github.com/RushabhMehta2005/post-go-res/memstore"


type NoOpWAL struct {}

func NewNoOpWAL() (*NoOpWAL, error) {
	w := new(NoOpWAL)
	var err error
	return w, err
}

func (w *NoOpWAL) Log(entry LogEntry) error {
	return nil
}

func (w *NoOpWAL) ReBuild(store store.InMemStore) error {
	return nil
}
