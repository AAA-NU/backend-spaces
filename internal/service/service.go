package service

import (
	"context"
	"log/slog"

	"github.com/aaanu/backend-spaces/internal/domain/models"
)

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

func (s *SpacesService) Coworking(ctx context.Context, id string) (*models.CoworkingModel, error) {
	const op = "service.Coworking"
	log := s.log.With("op", op)

	coworking, err := s.storage.Coworking(ctx, id)
	if err != nil {
		log.Error("failed to get coworking", "error", err)
		return nil, err
	}

	log.Info("got coworking", "id", id)
	return coworking, nil
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

	if err := s.storage.AddCoworkingBookedTime(ctx, coworkingID, newTime); err != nil {
		log.Error("failed to add coworking booked time", "error", err)
		return err
	}

	log.Info("added coworking booked time", "coworking_id", coworkingID, "time", newTime)
	return nil
}
