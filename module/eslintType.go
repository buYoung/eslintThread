package module

import "github.com/radovskyb/watcher"

type EslintInstance struct {
	EslintPath   string           `json:"path"`
	IgnoreList   []string         `json:"ignore_list"`
	FileWatcher  *watcher.Watcher `json:"_"`
	isInitialize bool
}
