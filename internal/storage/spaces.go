package storage

import (
	"context"
	"errors"
	"slices"
	"time"

	"github.com/aaanu/backend-spaces/internal/domain/models"
	"github.com/lib/pq"
)

func (s *Storage) Room(ctx context.Context, id string) (*models.RoomModel, error) {
	var room models.RoomModel
	if err := s.db.WithContext(ctx).Where("id = ?", id).First(&room).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

func (s *Storage) Coworking(ctx context.Context, id string) (*models.CoworkingModel, error) {
	var coworking models.CoworkingModel
	if err := s.db.WithContext(ctx).Where("id = ?", id).First(&coworking).Error; err != nil {
		return nil, err
	}
	return &coworking, nil
}

func (s *Storage) Rooms(ctx context.Context) ([]models.RoomModel, error) {
	var rooms []models.RoomModel
	if err := s.db.WithContext(ctx).Find(&rooms).Error; err != nil {
		return nil, err
	}
	return rooms, nil
}

func (s *Storage) Coworkings(ctx context.Context) ([]models.CoworkingModel, error) {
	var coworkings []models.CoworkingModel
	if err := s.db.WithContext(ctx).Find(&coworkings).Error; err != nil {
		return nil, err
	}
	return coworkings, nil
}

func (s *Storage) UpdateRoomBooking(ctx context.Context, room *models.RoomModel) error {
	return s.db.WithContext(ctx).Model(&models.RoomModel{}).Where("id = ?", room.ID).Save(room).Error
}

func (s *Storage) AddCoworkingBookedTime(ctx context.Context, coworkingID string, newTime string) error {
	var current models.CoworkingModel
	if err := s.db.First(&current, "id = ?", coworkingID).Error; err != nil {
		return err
	}

	now := time.Now()
	validBookings := make(pq.StringArray, 0, len(current.BookedTime))

	for _, timeStr := range current.BookedTime {
		bookedTime, err := time.Parse("2006-01-02 15:04", timeStr)
		if err != nil || bookedTime.After(now) {
			validBookings = append(validBookings, timeStr)
		}
	}

	if slices.Contains(validBookings, newTime) {
		return errors.New("coworking is already booked")
	}

	validBookings = append(validBookings, newTime)
	current.BookedTime = validBookings
	return s.db.Save(&current).Error
}
