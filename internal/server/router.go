package server

import (
	"github.com/aaanu/backend-spaces/internal/domain/models"
	"github.com/aaanu/backend-spaces/internal/domain/requests"
	"github.com/gin-gonic/gin"
)

type Router struct {
	router  *gin.RouterGroup
	service SpacesService
}

func Register(router *gin.RouterGroup, service SpacesService) {
	r := Router{router: router, service: service}
	r.init()
}

func (r *Router) GetRoom(ctx *gin.Context) {
	roomID := ctx.Param("id")

	room, err := r.service.Room(ctx, roomID)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	ctx.JSON(200, room)
}

func (r *Router) GetRooms(ctx *gin.Context) {
	rooms, err := r.service.Rooms(ctx)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	ctx.JSON(200, rooms)
}

func (r *Router) GetCoworking(ctx *gin.Context) {
	coworkingID := ctx.Param("id")
	date := ctx.Query("date")

	if date == "" {
		ctx.JSON(400, gin.H{"status": "error", "message": "query param \"date\" is required"})
		return
	}

	coworking, err := r.service.Coworking(ctx, coworkingID, date)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	ctx.JSON(200, coworking)
}

func (r *Router) GetCoworkings(ctx *gin.Context) {
	coworkings, err := r.service.Coworkings(ctx)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	ctx.JSON(200, coworkings)
}

func (r *Router) UpdateRoomBooking(ctx *gin.Context) {
	var room models.RoomModel
	if err := ctx.ShouldBindJSON(&room); err != nil {
		HandleError(ctx, err)
		return
	}

	if err := r.service.UpdateRoomBooking(ctx, &room); err != nil {
		HandleError(ctx, err)
		return
	}

	ctx.JSON(200, room)
}

func (r *Router) AddCoworkingBookedTime(ctx *gin.Context) {
	var newTime requests.AddBookingTime

	if err := ctx.ShouldBindJSON(&newTime); err != nil {
		HandleError(ctx, err)
		return
	}

	if err := r.service.AddCoworkingBookedTime(ctx, ctx.Param("id"), newTime.Time); err != nil {
		HandleError(ctx, err)
		return
	}

	ctx.JSON(200, newTime)
}

func (r *Router) init() {
	roomsGroup := r.router.Group("/rooms")

	roomsGroup.GET("/:id", r.GetRoom)
	roomsGroup.GET("/", r.GetRooms)
	roomsGroup.PUT("/", r.UpdateRoomBooking)

	coworkingsGroup := r.router.Group("/coworkings")

	coworkingsGroup.GET("/:id", r.GetCoworking)
	coworkingsGroup.GET("/", r.GetCoworkings)
	coworkingsGroup.POST("/:id", r.AddCoworkingBookedTime)
}
