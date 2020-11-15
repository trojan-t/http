package banners

import (
	"context"
	"errors"
	"sync"
)

// Banner is banner struct
type Banner struct {
	ID      int64
	Title   string
	Content string
	Button  string
	Link    string
	Image   string
}

// Service is service struct
type Service struct {
	mu    sync.RWMutex
	items []*Banner
}

// NewService is function
func NewService() *Service {
	return &Service{items: make([]*Banner, 0)}
}

// GetAll is CRUD method
func (s *Service) GetAll(ctx context.Context) ([]*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.items, nil
}

// GetByID is CRUD method
func (s *Service) GetByID(ctx context.Context, id int64) (*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, banner := range s.items {
		if banner.ID == id {
			return banner, nil
		}
	}
	return nil, errors.New("banner by id not found")
}

var starID int64 = 0

// Save is CRUD method
func (s *Service) Save(ctx context.Context, item *Banner) (*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if item.ID == 0 {
		starID++
		item.ID = starID
		s.items = append(s.items, item)
		return item, nil
	}
	for i, banner := range s.items {
		if banner.ID == item.ID {
			s.items[i] = item
			return item, nil
		}
	}

	return nil, errors.New("banner save error")
}

// RemoveByID is CRUD method
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
