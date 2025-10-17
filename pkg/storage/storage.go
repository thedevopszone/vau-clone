package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

// Storage represents the storage backend interface
type Storage interface {
	Get(key string) ([]byte, error)
	Put(key string, value []byte) error
	Delete(key string) error
	List(prefix string) ([]string, error)
}

// FileStorage implements file-based storage
type FileStorage struct {
	basePath string
	mu       sync.RWMutex
	data     map[string][]byte
}

// NewFileStorage creates a new file storage backend
func NewFileStorage(basePath string) (*FileStorage, error) {
	if err := os.MkdirAll(basePath, 0700); err != nil {
		return nil, err
	}

	fs := &FileStorage{
		basePath: basePath,
		data:     make(map[string][]byte),
	}

	// Load existing data
	if err := fs.load(); err != nil {
		return nil, err
	}

	return fs, nil
}

// Get retrieves a value by key
func (fs *FileStorage) Get(key string) ([]byte, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	value, exists := fs.data[key]
	if !exists {
		return nil, errors.New("key not found")
	}

	return value, nil
}

// Put stores a value by key
func (fs *FileStorage) Put(key string, value []byte) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	fs.data[key] = value
	return fs.persist()
}

// Delete removes a value by key
func (fs *FileStorage) Delete(key string) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	if _, exists := fs.data[key]; !exists {
		return errors.New("key not found")
	}

	delete(fs.data, key)
	return fs.persist()
}

// List returns all keys with the given prefix
func (fs *FileStorage) List(prefix string) ([]string, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	var keys []string
	for key := range fs.data {
		if len(prefix) == 0 || len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			keys = append(keys, key)
		}
	}

	return keys, nil
}

// persist saves the data to disk
func (fs *FileStorage) persist() error {
	dataFile := filepath.Join(fs.basePath, "vault.db")
	data, err := json.MarshalIndent(fs.data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(dataFile, data, 0600)
}

// load loads data from disk
func (fs *FileStorage) load() error {
	dataFile := filepath.Join(fs.basePath, "vault.db")

	// If file doesn't exist, that's okay
	if _, err := os.Stat(dataFile); os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(dataFile)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &fs.data)
}
