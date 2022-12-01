package module

import (
	"errors"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/radovskyb/watcher"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

func (e *FileWatcherInstance) Init() error {
	e.isInitialize = true
	e.FileWatcher.FilterOps(watcher.Write)
	err := e.ReadIgnoreFile()
	if err != nil {
		return err
	}
	return nil
}

func (e *FileWatcherInstance) Start() error {
	if !e.isInitialize {
		return errors.New("초기화 필요")
	}
	go func() {
		for {
			select {
			case event := <-e.FileWatcher.Event:
				if event.IsDir() {
					continue
				}
				folderPath := event.Path
				isIgnoreFile := e.isIgnoreFile(folderPath)
				if isIgnoreFile {
					continue
				}
				e.DB.DB.Update(func(txn *badger.Txn) error {
					e.dbUpdateValue(folderPath, false)
					return nil
				})
			case err := <-e.FileWatcher.Error:
				log.Println(err)
			case <-e.FileWatcher.Closed:
				return
			}
		}
	}()

	if err := e.FileWatcher.AddRecursive(filepath.Clean(filepath.Join(e.EslintPath))); err != nil {
		log.Println("error read Folder ", err)
		return err
	}
	//
	var pathList []string
	for _, info := range e.FileWatcher.WatchedFiles() {
		if info.IsDir() {
			continue
		}
		value := reflect.ValueOf(info).Elem()
		if value.Kind() == reflect.Struct {
			vaildPath := value.FieldByName("path")
			if vaildPath.IsValid() {
				pathList = append(pathList, filepath.Clean(vaildPath.String()))
			}
		}
	}
	e.sortFilename(pathList)
	sort.Sort(pathSorter(pathList))
	e.sortFilePathDeps(pathList)
	for _, s := range pathList {
		e.dbUpdateValue(s, true)
	}
	if err := e.FileWatcher.Start(time.Millisecond * 10); err != nil {
		return err
	}
	return nil
}
func (e *FileWatcherInstance) dbUpdateValue(key string, isInit bool) {
	err := e.DB.DB.Update(func(txn *badger.Txn) error {
		defaultValue := 0
		var valCopy []byte
		get, err := txn.Get([]byte(key))
		if err != badger.ErrKeyNotFound {
			err = get.Value(func(val []byte) error {
				valCopy = append([]byte{}, val...)
				return nil
			})
			if err != nil {
				if isInit {
					defaultValue, err = strconv.Atoi(string(valCopy))
					if err != nil {
						defaultValue = 0
					}
					defaultValue += 1
				}
			}
		}

		return txn.Set([]byte(key), []byte(strconv.Itoa(defaultValue)))
	})
	if err != nil {
		log.Println("db update Error")
	}

}
func (e *FileWatcherInstance) sortFilePathDeps(filePathList []string) {
	sort.SliceStable(filePathList, func(i, j int) bool {
		a := strings.Count(filePathList[i], string(os.PathSeparator))
		b := strings.Count(filePathList[j], string(os.PathSeparator))
		return a < b
	})
}
func (e *FileWatcherInstance) sortFilename(filePathList []string) {
	sort.SliceStable(filePathList, func(i, j int) bool {
		a := filepath.Base(filePathList[i])
		b := filepath.Base(filePathList[j])
		return a < b
	})
}

func (e *FileWatcherInstance) isEslintIgnoreFile() (string, error) {
	ignorePath := filepath.Clean(filepath.Join(e.EslintPath, ".eslintignore"))
	_, err := os.Stat(ignorePath)
	if err != nil {
		return "", errors.New("notFound")
	}
	return ignorePath, nil
}
func (e *FileWatcherInstance) ReadIgnoreFile() (err error) {
	defer func() {
		s := recover()
		if s != nil {
			err = errors.New(fmt.Sprintf("readIgnoreFileError %v", s))
		}
	}()
	ignorePathData, err := e.isEslintIgnoreFile()
	if err != nil {
		return err
	}
	eslintignoreData, err := os.ReadFile(ignorePathData)
	if err != nil {
		log.Println("파일 읽다가 오류생김", err)
	}

	eslintignoreDataSplitList := strings.Split(string(eslintignoreData), "\n")

	var ignorePathList []string
	for _, s := range eslintignoreDataSplitList {
		isDir, err := os.Stat(filepath.Clean(filepath.Join(e.EslintPath, s)))
		if err == nil && isDir.IsDir() {
			ignorePath := filepath.Clean(filepath.Join(e.EslintPath, s))
			ignorePathList = append(ignorePathList, ignorePath)
		}
		e.IgnoreList = append(e.IgnoreList, strings.Trim(s, " "))
	}

	err = e.FileWatcher.Ignore(ignorePathList...)
	if err != nil {
		log.Println("error ignore regist", err, ignorePathList)
	}
	return nil
}

func (e *FileWatcherInstance) isIgnoreFile(dir string) bool {
	for _, s := range e.IgnoreList {
		if strings.Contains(dir, s) {
			return true
		}
	}
	return false

}
