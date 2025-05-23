package db

import "sync"

type Manager struct {
	wg    sync.WaitGroup
	mutex sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) Create(data []byte) error {
	return Create(data)
}

func (m *Manager) Read(index uint64) ([]byte, error) {
	return Read(index)
}

func (m *Manager) Update(index uint64, data []byte) error {
	return Update(index, data, m)
}

func (m *Manager) Delete(index uint64) {
	Delete(index)
}
