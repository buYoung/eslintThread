package main

import (
	"eslintThread/data"
	"eslintThread/module"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/gopherjs/gopherjs/js"
	"github.com/radovskyb/watcher"
	"log"
	"time"
)

var (
	Db data.BadgerInstance
)

func main() {

	/////TODO https://gist.github.com/rushilgupta/228dfdf379121cb9426d5e90d34c5b96
	Db, err := data.BadgerIitialize()
	if err != nil {
		log.Println(err)
	}

	eslintTracker := &module.FileWatcherInstance{
		EslintPath:  "Z:\\innerview\\moye\\app",
		FileWatcher: watcher.New(),
		DB:          Db,
	}
	err = eslintTracker.Init()
	if err != nil {
		log.Println("Init", err)
	}
	err = eslintTracker.Start()
	if err != nil {
		log.Println("Init", err)
	}

	var filelist []string
	go func() {
		for {
			time.Sleep(time.Second * 5)
			eslintTracker.DB.DB.View(func(txn *badger.Txn) error {
				opt := badger.IteratorOptions{
					PrefetchSize: 10,
				}
				it := txn.NewIterator(opt)
				log.Println(it)
				if !it.Valid() {
					return nil
				}
				log.Println(it)
				defer it.Close()
				for it.Rewind(); it.Valid(); it.Next() {
					item := it.Item()
					k := item.Key()
					err := item.Value(func(v []byte) error {
						fmt.Printf("key=%s, value=%s\n", k, v)
						return nil
					})
					if err != nil {
						return err
					}
				}
				return nil
			})
		}
	}()

	js.Module.Get("exports").Set("fileList", filelist)
}
