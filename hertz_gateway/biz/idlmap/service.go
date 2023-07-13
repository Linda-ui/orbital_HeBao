package idlmap

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/fsnotify/fsnotify"
)

type Repository interface {
	GetClient(svcName string) (cli genericclient.Client, ok bool)
	AddService(idlPath string, opts ...client.Option) error
	DeleteService(svcName string)
}

// manager implements the entity.MapManager interface.
type manager struct {
	repo Repository
}

func NewManager(r Repository) *manager {
	return &manager{r}
}

func (m *manager) GetClient(svcName string) (genericclient.Client, bool) {
	return m.repo.GetClient(svcName)
}

func (m *manager) AddService(idlPath string, opts ...client.Option) error {
	return m.repo.AddService(idlPath, opts...)
}

func (m *manager) DeleteService(svcName string) {
	m.repo.DeleteService(svcName)
}

func (m *manager) AddAllServices(idlRootPath string, opts ...client.Option) {
	entries, err := os.ReadDir(idlRootPath)
	if err != nil {
		log.Fatalf("scanning idl file directory failed: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			m.AddAllServices(idlRootPath+"/"+entry.Name(), opts...)
		} else {
			m.repo.AddService(idlRootPath+"/"+entry.Name(), opts...)
		}
	}
}

func (m *manager) DynamicUpdate(idlRootPath string, opts ...client.Option) {
	dirPathSet := make(map[string]struct{})

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
				path := event.Name
				switch {
				case event.Has(fsnotify.Create) || event.Has(fsnotify.Write):
					if isDirectory(path) {
						hlog.Infof("Directory " + event.String())
						watcher.Add(path)
						dirPathSet[path] = struct{}{}
					} else {
						hlog.Infof("File " + event.String())
						m.repo.AddService(path, opts...)
					}
				case event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename):
					_, ok := dirPathSet[path]
					if ok {
						hlog.Infof("Directory " + event.String())
						delete(dirPathSet, path)
					} else {
						hlog.Infof("File " + event.String())
						svcName := strings.ReplaceAll(filepath.Base(path), ".thrift", "")
						m.repo.DeleteService(svcName)
					}
				}

			case err := <-watcher.Errors:
				hlog.Infof("Error: %v", err)
			}
		}
	}()

	filepath.Walk(idlRootPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			watcher.Add(path)
			dirPathSet[path] = struct{}{}
		}
		return nil
	})

	<-done
}

func isDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Print(err)
	}
	return fileInfo.IsDir()
}
