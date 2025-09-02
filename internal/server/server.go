package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aaanu/backend-spaces/internal/config"
	"github.com/aaanu/backend-spaces/internal/domain/models"
	"github.com/gin-gonic/gin"
)

type SpacesService interface {
	Room(ctx context.Context, id string) (*models.RoomModel, error)
	Rooms(ctx context.Context) ([]models.RoomModel, error)
	Coworking(ctx context.Context, id string) (*models.CoworkingModel, error)
	Coworkings(ctx context.Context) ([]models.CoworkingModel, error)
	UpdateRoomBooking(ctx context.Context, room *models.RoomModel) error
	AddCoworkingBookedTime(ctx context.Context, coworkingID string, newTime string) error
}

type Server struct {
	server *http.Server
	engine *gin.Engine
}

func New(
	service SpacesService,
) *Server {
	cfg := config.Config().Server
	engine := gin.Default()

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: engine,
	}

	group := engine.Group("/api")
	group.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	return &Server{
		server: httpServer,
		engine: engine,
	}
}

func (s *Server) Start() {
	if err := s.server.ListenAndServe(); err != nil {
		panic(err)
	}
}

func (s *Server) GracefulStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		panic(err)
	}
}
