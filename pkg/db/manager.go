package db

import (
	"bufio"
	"encoding/json"
	"go_text_task/pkg/db/linked_list"
	"go_text_task/pkg/files"
	"os"
	"sync"
)

type Storage struct {
	ID       uint64
	IndexMap map[uint64]DataLocation
}
type Manager struct {
	Storage
	DL    *linked_list.DoubleLinkedList
	wg    *sync.WaitGroup
	mtx   *sync.RWMutex
	tasks chan func()
}

func NewManager() *Manager {

	tasks := make(chan func(), 3)
	go func(tasks chan func()) {
		for nextTask := range tasks {
			nextTask()
		}
	}(tasks)

	manager := Manager{}
	manager.wg = new(sync.WaitGroup)
	manager.mtx = new(sync.RWMutex)
	manager.DL = linked_list.NewDoubleLinkedList()
	manager.tasks = tasks

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

func (m *Manager) storeIndexes() {

	indexFile, creationErr := os.Create(createPath(IndexFileName))
	if creationErr != nil {
		panic(creationErr)
	}
	defer func(indexFile *os.File) {
		closeErr := indexFile.Close()
		if closeErr != nil {
			panic(closeErr)
		}
	}(indexFile)

	writer := bufio.NewWriter(indexFile)
	encodingErr := json.NewEncoder(writer).Encode(m.Storage)
	if encodingErr != nil {
		panic(encodingErr)
	}

	flushErr := writer.Flush()
	if flushErr != nil {
		panic(flushErr)
	}

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
