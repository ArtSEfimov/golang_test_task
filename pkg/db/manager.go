package db

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"
	"sync"
)

var IndexFilePath = createPath(IndexFileName)

type Storage struct {
	ID       uint64
	IndexMap map[uint64]DataLocation
}
type Manager struct {
	Storage
	wg    sync.WaitGroup
	mutex sync.RWMutex
}

func (m *Manager) GetID() uint64 {
	return m.ID
}

func NewManager() *Manager {

	manager := Manager{}

	if !isFileExists(IndexFilePath) {
		creationErr := createFile(IndexFileName)
		if creationErr != nil {
			panic(creationErr)
		}
		manager.IndexMap = make(map[uint64]DataLocation)
		manager.ID = 0
	} else {
		indexFile, openErr := os.Open(IndexFilePath)
		if openErr != nil {
			panic(openErr)
		}
		defer func(indexFile *os.File) {
			closeErr := indexFile.Close()
			if closeErr != nil {
				panic(closeErr)
			}
		}(indexFile)

		fileReader := bufio.NewReader(indexFile)
		decodeErr := json.NewDecoder(fileReader).Decode(&manager.Storage)
		if decodeErr != nil {
			if decodeErr == io.EOF {
				log.Printf("database index file is empty: %v", decodeErr)
			} else {
				panic(decodeErr)
			}
		}
	}

	return &manager
}

func (m *Manager) storeIndexes() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	indexFile, creationErr := os.Create(IndexFilePath)
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

	err := writer.Flush()
	if err != nil {
		panic(err)
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
