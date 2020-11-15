package banners

import (
	"context"
	"errors"
	"sync"
)

// Service presents struct
type Service struct {
	mu    sync.RWMutex
	items []*Banner
}

// NewService creates new service
func NewService() *Service {
	return &Service{items: make([]*Banner, 0)}
}

// Banner presents banner :)
type Banner struct {
	ID      int64
	Title   string
	Content string
	Button  string
	Link    string
}

// All returns all banners
func (s *Service) All(ctx context.Context) ([]*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.items, nil
}

// Save save/updates banner
func (s *Service) Save(ctx context.Context, item *Banner) (*Banner, error) {
	var i int64 = 0
	s.mu.RLock()
	defer s.mu.RUnlock()
	if item.ID == 0 {
		i++
		item.ID = i
		s.items = append(s.items, item)
		return item, nil
	}

	for index, banner := range s.items {
		if banner.ID == item.ID {
			s.items[index] = item
			return item, nil
		}
	}
	return nil, errors.New("item not found")
}

// ByID return banner by ID
func (s *Service) ByID(ctx context.Context, id int64) (*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, banner := range s.items {
		if banner.ID == id {
			return banner, nil
		}
	}
	return nil, errors.New("item not found")
}

// RemoveByID removes banner by ID
func (s *Service) RemoveByID(ctx context.Context, id int64) (*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for index, banner := range s.items {
		if banner.ID == id {
			s.items = append(s.items[:index], s.items[index+1:]...)
			return banner, nil
		}
	}
	return nil, errors.New("item not found")
}
