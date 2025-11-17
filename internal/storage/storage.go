package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"links_project/internal/models"
	"os"
	"sync"
)


type Storage interface {
	Load() error
	Save() error
	NextID() int
	SaveBatch(b *models.Batch)
	GetBatch(id int) (*models.Batch, bool)
}

type storage struct {
	Path string
	mu sync.RWMutex
	Batches map[int]*models.Batch `json:"batches"`
	LastID int `json:"last_id"`
}

func NewStorage(path string) Storage{
	return &storage{
		Path: path,
		Batches: make(map[int]*models.Batch),
	}
}

func (s *storage) Load() error{
	s.mu.Lock()
	defer s.mu.Unlock()
	data, err := os.ReadFile(s.Path)
	if err != nil{
		if errors.Is(err, os.ErrNotExist){
			return nil
		}
		return fmt.Errorf("failed to load data: %w", err)
	}

	if len(data) == 0 {
		return nil
	}
	err = json.Unmarshal(data, s)
	if err != nil {
		return fmt.Errorf("failed to unmarshal data: %w", err)
	}
	return nil
}

func (s *storage) Save() error{
	s.mu.RLock()
	defer s.mu.RUnlock()
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil{
		return fmt.Errorf("failed to marshal data: %w", err)
	}
	err = os.WriteFile(s.Path, data, 0644)
	if err != nil{
		return fmt.Errorf("failed to save data: %w", err)
	}
	return nil
}

func (s *storage) NextID() int{
	s.mu.Lock()
	defer s.mu.Unlock()
	s.LastID++
	return s.LastID
}

func (s *storage) SaveBatch(b *models.Batch){
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Batches[b.ID] = b
}

func (s *storage) GetBatch(id int) (*models.Batch, bool){
	s.mu.RLock()
	defer s.mu.RUnlock()
	data, ok := s.Batches[id]
	return data, ok
}