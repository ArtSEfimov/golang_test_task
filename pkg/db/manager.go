package db

import (
	"bufio"
	"encoding/json"
	"go_text_task/pkg/db/config"
	"go_text_task/pkg/db/linked_list"
	"go_text_task/pkg/files"
	"os"
	"sync"
)

var (
	IndexFileName            string
	OriginalDataBaseFileName string
	MaxFileSegmentSize       uint64
	DataBaseDir              string
)

type Manager struct {
	Storage
	DL    *linked_list.DoubleLinkedList
	wg    *sync.WaitGroup
	mtx   *sync.RWMutex
	Tasks chan func()
	Done  chan struct{}
}

func NewManager(config *config.Config) *Manager {

	IndexFileName = config.GetIndexFileName()
	OriginalDataBaseFileName = config.GetOriginalDBFileName()
	MaxFileSegmentSize = config.GetMaxFileSegmentSize()
	DataBaseDir = config.GetDBDir()

	tasks := make(chan func())
	done := make(chan struct{})
	go func() {
		for nextTask := range tasks {
			nextTask()
		}
		done <- struct{}{}
	}()

	manager := Manager{}
	manager.wg = new(sync.WaitGroup)
	manager.mtx = new(sync.RWMutex)
	manager.DL = linked_list.NewDoubleLinkedList()
	manager.Tasks = tasks
	manager.Done = done

	if !files.IsFileExists(createPath(IndexFileName)) {
		manager.IndexMap = make(map[uint64]DataLocation)
		manager.ID = 0
	} else {
		indexFile, openErr := os.Open(createPath(IndexFileName))
		if openErr != nil {
			panic(openErr)
		}
		defer func(file *os.File) {
			closeErr := file.Close()
			if closeErr != nil {
				panic(closeErr)
			}
		}(indexFile)

		fileReader := bufio.NewReader(indexFile)
		decodeErr := json.NewDecoder(fileReader).Decode(&manager.Storage)
		if decodeErr != nil {
			panic(decodeErr)
		}
	}

	return &manager
}

func (m *Manager) GetID() uint64 {
	return m.ID
}

func (m *Manager) Create(data []byte) error {
	return Create(data, m)
}

func (m *Manager) Read(index uint64) ([]byte, error) {
	return Read(index, m)
}

func (m *Manager) Update(index uint64, data []byte) error {
	return Update(index, data, m)
}

func (m *Manager) Delete(index uint64) error {
	return Delete(index, m)
}
