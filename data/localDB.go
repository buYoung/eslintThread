package data

import (
	"github.com/dgraph-io/badger/v3"
)

type FileInfo struct {
	Path        string `json:"path,omitempty"`
	UpdateCount int    `json:"update_count"`
}

type BadgerInstance struct {
	DB *badger.DB
}

func BadgerIitialize() (*BadgerInstance, error) {
	option := badger.DefaultOptions("badger.db")
	option.WithInMemory(true)
	option.IndexCacheSize = 100 << 20
	db, err := badger.Open(option)
	if err != nil {
		return nil, err
	}
	B := &BadgerInstance{DB: db}

	return B, nil
}
