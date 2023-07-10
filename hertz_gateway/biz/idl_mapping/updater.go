package idl_mapping

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/client"
	"github.com/fsnotify/fsnotify"
)

func WatchAndUpdate(m IMap, opts ...client.Option) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer watcher.Close()

	done := make(chan struct{})
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Has(fsnotify.Create) {
					if !isDirectory(event.Name) {
						hlog.Infof("File created: %v", event.Name)
						m.Add(event.Name, opts...)
					} else {
						hlog.Infof("Directory created: %v", event.Name)
					}
				}
				if event.Has(fsnotify.Remove) {
					if !isDirectory(event.Name) {
						hlog.Infof("File removed: %v", event.Name)
						m.Delete(event.Name)
					} else {
						hlog.Infof("Directory created: %v", event.Name)
					}
				}
				if event.Has(fsnotify.Rename) {
					hlog.Infof("File renamed: %v", event.Name)
				}
			case err := <-watcher.Errors:
				log.Printf("Error: %v", err)
			}
		}
	}()

	filepath.Walk("./idl", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			watcher.Add(path)
		}
		return nil
	})

	<-done
}

func isDirectory(path string) bool {
	// fileInfo, err := os.Stat(path)
	// if err != nil {
	// 	log.Print(err)
	// }

	// return fileInfo.IsDir()
	return !strings.HasSuffix(filepath.Base(path), ".thrift")
}
