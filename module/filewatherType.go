package module

import (
	"eslintThread/data"
	"github.com/radovskyb/watcher"
)

type FileWatcherInstance struct {
	EslintPath   string           `json:"path"`
	IgnoreList   []string         `json:"ignore_list"`
	FileWatcher  *watcher.Watcher `json:"_"`
	isInitialize bool
	DB           *data.BadgerInstance
}
