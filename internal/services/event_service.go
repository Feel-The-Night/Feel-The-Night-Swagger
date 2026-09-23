package services

import (
	"errors"

	"github.com/kisalto/Feel-The-Night-Swagger/internal/models"
	"gorm.io/gorm"
)

type EventService struct {
	db *gorm.DB
}

func NewEventService(db *gorm.DB) *EventService {
	return &EventService{db: db}
}

func (s *EventService) CreateEvent(event *models.Event) error {
	return s.db.Create(event).Error
}

func (s *EventService) GetEventById(id uint) (*models.Event, error) {
	var event models.Event

	if err := s.db.First(&event, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &event, nil
}

func (s *EventService) DeleteEventById(id uint) error {
	result := s.db.Delete(&models.Event{EventID: id})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (s *EventService) UpdateEventById(id uint, event *models.Event) (*models.Event, error) {
	existingEvent, err := s.GetEventById(id)
	if err != nil {
		return nil, err
	}

	updates := make(map[string]any)

	if event.Title != "" {
		updates["title"] = event.Title
	}
	if event.BannerURL != "" {
		updates["banner_url"] = event.BannerURL
	}
	if !event.Day.IsZero() {
		updates["day"] = event.Day
	}
	if event.Description != "" {
		updates["description"] = event.Description
	}

	if len(updates) == 0 {
		return existingEvent, nil
	}

	if err := s.db.Model(existingEvent).Updates(updates).Error; err != nil {
		return nil, err
	}

	return existingEvent, nil
}
