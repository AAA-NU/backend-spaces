package service

import (
	"context"
	"errors"
	"log/slog"
	"slices"
	"time"

	"github.com/aaanu/backend-spaces/internal/domain/models"
	"github.com/aaanu/backend-spaces/internal/domain/responses"
)

var coworkingOpeningHours = []string{
	"08:00",
	"09:00",
	"10:00",
	"11:00",
	"12:00",
	"13:00",
	"14:00",
	"15:00",
	"16:00",
	"17:00",
	"18:00",
	"19:00",
	"20:00",
	"21:00",
	"22:00",
}

type SpacesStorage interface {
	Room(ctx context.Context, id string) (*models.RoomModel, error)
	Rooms(ctx context.Context) ([]models.RoomModel, error)

	Coworking(ctx context.Context, id string) (*models.CoworkingModel, error)
	Coworkings(ctx context.Context) ([]models.CoworkingModel, error)

	UpdateRoomBooking(ctx context.Context, room *models.RoomModel) error
	AddCoworkingBookedTime(ctx context.Context, coworkingID string, newTime string) error
}

type SpacesService struct {
	log     *slog.Logger
	storage SpacesStorage
}

func New(log *slog.Logger, storage SpacesStorage) *SpacesService {
	log = log.With("service", "spaces")
	return &SpacesService{log: log, storage: storage}
}

func (s *SpacesService) Room(ctx context.Context, id string) (*models.RoomModel, error) {
	const op = "service.Room"
	log := s.log.With("op", op)

	room, err := s.storage.Room(ctx, id)
	if err != nil {
		log.Error("failed to get room", "error", err)
		return nil, err
	}

	log.Info("got room", "id", id)
	return room, nil
}

func (s *SpacesService) Rooms(ctx context.Context) ([]models.RoomModel, error) {
	const op = "service.Rooms"
	log := s.log.With("op", op)

	rooms, err := s.storage.Rooms(ctx)
	if err != nil {
		log.Error("failed to get rooms", "error", err)
		return nil, err
	}

	log.Info("got rooms")
	return rooms, nil
}

func (s *SpacesService) Coworking(ctx context.Context, id string, date string) (*responses.CoworkingMetaResponse, error) {
	const op = "service.Coworking"
	log := s.log.With("op", op)

	if _, err := time.Parse("2006-01-02", date); err != nil {
		log.Error("failed to parse date", "error", err)
		return nil, err
	}

	coworking, err := s.storage.Coworking(ctx, id)
	if err != nil {
		log.Error("failed to get coworking", "error", err)
		return nil, err
	}

	// coworking works from 8:00 till 22:00. Loop through booked time of coworking and find if there are any bookings
	// for the given date. Date in bookings presented as '2006-01-02 15:04'. you have to put all available times in response excluding booked times
	var bookedTimes []string

	for _, timeStr := range coworking.BookedTime {
		bookedTime, err := time.Parse("2006-01-02 15:04", timeStr)
		if err != nil {
			log.Error("failed to parse time", "error", err)
			return nil, err
		}

		if bookedTime.Format("2006-01-02") == date {
			bookedTimes = append(bookedTimes, bookedTime.Format("15:04"))
		}
	}

	var availableTimes []string
	for _, openingHour := range coworkingOpeningHours {
		if !slices.Contains(bookedTimes, openingHour) {
			availableTimes = append(availableTimes, openingHour)
		}
	}

	response := &responses.CoworkingMetaResponse{
		ID:            coworking.ID,
		AvailableTime: availableTimes,
	}

	log.Info("got coworking", "id", id)
	return response, nil
}

func (s *SpacesService) Coworkings(ctx context.Context) ([]models.CoworkingModel, error) {
	const op = "service.Coworkings"
	log := s.log.With("op", op)

	coworkings, err := s.storage.Coworkings(ctx)
	if err != nil {
		log.Error("failed to get coworkings", "error", err)
		return nil, err
	}

	log.Info("got coworkings")
	return coworkings, nil
}

func (s *SpacesService) UpdateRoomBooking(ctx context.Context, room *models.RoomModel) error {
	const op = "service.UpdateRoomBooking"
	log := s.log.With("op", op)

	if err := s.storage.UpdateRoomBooking(ctx, room); err != nil {
		log.Error("failed to update room booking", "error", err)
		return err
	}

	log.Info("updated room booking", "id", room.ID)
	return nil
}

func (s *SpacesService) AddCoworkingBookedTime(ctx context.Context, coworkingID string, newTime string) error {
	const op = "service.AddCoworkingBookedTime"
	log := s.log.With("op", op)

	if len(newTime) <= 10 {
		log.Error("failed to add coworking booked time", "error", "invalid time")
		return errors.New("invalid time")
	}

	if err := s.storage.AddCoworkingBookedTime(ctx, coworkingID, newTime); err != nil {
		log.Error("failed to add coworking booked time", "error", err)
		return err
	}

	log.Info("added coworking booked time", "coworking_id", coworkingID, "time", newTime)
	return nil
}
