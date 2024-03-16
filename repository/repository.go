package repository

import (
	"path/filepath"
	"sync"

	cmap "github.com/orcaman/concurrent-map"
)

// QueueRepository represents a repository of mesage queues
type QueueRepository struct {
	sync.Mutex
	// cmap.ConcurrentMap is a thread-safe concurrent map.
	// It provides a high-performance solution to this by sharding the map with minimal time spent waiting for locks
	storage  cmap.ConcurrentMap
	DataPath string
}

// NewRepository creates a new repository of queues
func NewRepository(dataDir string) (*QueueRepository, error) {
	dataPath, err := filepath.Abs(dataDir)
	if err != nil {
		return nil, err
	}
	repo := QueueRepository{storage: cmap.New(), DataPath: dataPath}
	return &repo, nil
}
