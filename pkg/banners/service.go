package banners

import (
	"context"
	"errors"
	"sync"
)

// STORAGE is path
const STORAGE = "./web/banners/"

// Banner is struct
type Banner struct {
	ID      int64
	Title   string
	Content string
	Button  string
	Link    string
}

// Service is struct
type Service struct {
	mu    sync.RWMutex
	items []*Banner
}

// NewService creates new service
func NewService() *Service {
	return &Service{items: make([]*Banner, 0)}
}

// All is method
func (s *Service) All(ctx context.Context) ([]*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.items, nil
}

// ByID is method
func (s *Service) ByID(ctx context.Context, id int64) (*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, banner := range s.items {
		if banner.ID == id {
			return banner, nil
		}
	}
	return nil, errors.New("banner not found")
}

// Save is method
func (s *Service) Save(ctx context.Context, item *Banner) (*Banner, error) {
	var ID int64
	s.mu.RLock()
	defer s.mu.RUnlock()
	if item.ID == 0 {
		ID++
		item.ID = ID
		s.items = append(s.items, item)
	} else if item.ID != 0 {
		s.items = append(s.items, item)
	}
	for i, banner := range s.items {
		if banner.ID == item.ID {
			s.items[i] = item
			return item, nil
		}
	}

	return nil, errors.New("banner save error")
}

// RemoveByID is method
func (s *Service) RemoveByID(ctx context.Context, id int64) (*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for i, banner := range s.items {
		if banner.ID == id {
			s.items = append(s.items[:i], s.items[i+1:]...)
			return banner, nil
		}
	}
	return nil, errors.New("banner remove by id not found")
}
