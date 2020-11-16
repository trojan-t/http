package banners

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"mime/multipart"
	"strconv"
	"sync"
)

// NotFound is
var NotFound = errors.New("Not Found")

// Service представляет собой сервис по управлению баннерами
type Service struct {
	mu     sync.RWMutex
	items  []*Banner
	NextID int64
}

// NewService создаёт сервис
func NewService() *Service {
	return &Service{items: make([]*Banner, 0)}
}

// Banner представляет собой баннер
type Banner struct {
	ID      int64
	Title   string
	Content string
	Button  string
	Link    string
	Image   string
}

// All возвращает все существующие баннеры
func (s *Service) All(ctx context.Context) ([]*Banner, error) {
	log.Println("мы в All")
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.items, nil
}

// Save сохраняет/обновляет баннер
func (s *Service) Save(ctx context.Context, item *Banner, image multipart.File) (*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if item.ID == 0 {
		s.NextID++
		item.ID = s.NextID

		if item.Image != "" {
			item.Image = strconv.Itoa(int(item.ID)) + "." + item.Image
			data, err := ioutil.ReadAll(image)
			if err != nil {
				return nil, NotFound
			}
			err = ioutil.WriteFile("./web/banners/"+item.Image, data, 0666)
			if err != nil {
				log.Print(err)
				return nil, NotFound
			}
		}
		s.items = append(s.items, item)
		return item, nil
	}

	for index, banner := range s.items {
		if banner.ID == item.ID {
			if item.Image != "" {

				item.Image = strconv.Itoa(int(item.ID)) + "." + item.Image
				data, err := ioutil.ReadAll(image)
				if err != nil {
					log.Print(err)
					return nil, err
				}
				err = ioutil.WriteFile("./web/banners/"+item.Image, data, 0666)
				if err != nil {
					log.Print(err)
					return nil, err
				}
			} else {
				item.Image = s.items[index].Image
			}
			s.items[index] = item
			return item, nil

		}

	}
	log.Println("мы возвращаем ошибку 8")
	return nil, NotFound
}

// ByID возвращает баннер по идентификатору
func (s *Service) ByID(ctx context.Context, id int64) (*Banner, error) {
	log.Println("мы в ByID")
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, banner := range s.items {
		if banner.ID == id {
			return banner, nil
		}
	}
	log.Println("мы в ByID2")
	return nil, NotFound
}

// RemoveByID удаляет баннер по идентификатору
func (s *Service) RemoveByID(ctx context.Context, id int64) (*Banner, error) {
	log.Println("мы в RemoveByID")
	s.mu.RLock()
	defer s.mu.RUnlock()
	for index, banner := range s.items {
		if banner.ID == id {
			s.items = append(s.items[:index], s.items[index+1:]...)
			return banner, nil
		}
	}
	return nil, NotFound
}
