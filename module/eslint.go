package module

import (
	"errors"
	"fmt"
	"github.com/radovskyb/watcher"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

func (e *EslintInstance) Init() error {
	e.isInitialize = true
	e.FileWatcher.FilterOps(watcher.Write)
	return nil
}

func (e *EslintInstance) Start() error {
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
				eventName := event.Op.String()
				fileName := event.Name()
				folderPath := event.Path
				isIgnoreFile := e.isIgnoreFile(folderPath)
				log.Println("get Event", isIgnoreFile, eventName, folderPath, fileName)
			case err := <-e.FileWatcher.Error:
				log.Println(err)
			case <-e.FileWatcher.Closed:
				return
			}
		}
	}()

	if err := e.FileWatcher.AddRecursive(path.Clean(path.Join(e.EslintPath))); err != nil {
		log.Println("error read Folder ", err)
		return err
	}
	if err := e.FileWatcher.Start(time.Millisecond * 10); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (e *EslintInstance) isEslintIgnoreFile() (string, error) {
	ignorePath := path.Clean(path.Join(e.EslintPath, ".eslintignore"))
	_, err := os.Stat(ignorePath)
	if err != nil {
		return "", errors.New("notFound")
	}
	return ignorePath, nil
}
func (e *EslintInstance) ReadIgnoreFile(fileDir string) (err error) {
	defer func() {
		s := recover()
		if s != nil {
			err = errors.New(fmt.Sprintf("readIgnoreFileError %v", s))
		}
	}()
	filePath, err := e.isEslintIgnoreFile()
	if err != nil {
		return err
	}
	eslintignoreData, err := os.ReadFile(filePath)
	if err != nil {
		log.Println("파일 읽다가 오류생김", err)
	}

	eslintignoreDataSplitList := strings.Split(string(eslintignoreData), "\n")

	for _, s := range eslintignoreDataSplitList {
		e.IgnoreList = append(e.IgnoreList, strings.Trim(s, " "))
	}
	return nil
}

func (e *EslintInstance) isIgnoreFile(dir string) bool {
	for _, s := range e.IgnoreList {
		if strings.Contains(dir, s) {
			return true
		}
	}
	return false

}
