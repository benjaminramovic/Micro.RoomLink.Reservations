package controllers

import (
	"encoding/json"
	"example/micro-roomlink-reservations/models"
	"example/micro-roomlink-reservations/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

type ReservationController struct {
	ReservationService services.ReservationService
	RabbitMQChannel *amqp.Channel
}

func NewReservationController(reservationService services.ReservationService, channel *amqp.Channel) *ReservationController {
	return &ReservationController{ReservationService: reservationService, RabbitMQChannel: channel}
}

func (rc *ReservationController) Create(ctx *gin.Context)  {
	/*_, err := rs.reservationsCollection.InsertOne(rs.ctx, reservation)
	if err != nil {
		return err
	}*/
	var reservation models.Reservation
	 err := ctx.ShouldBindJSON(&reservation)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	err = rc.ReservationService.Create(&reservation)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	// Provera da li je RabbitMQ kanal otvoren pre slanja poruke
	if rc.RabbitMQChannel == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "RabbitMQ channel is not initialized"})
		return
	}

	// Publish message to RabbitMQ
	message := map[string]interface{}{
		"roomId": reservation.RoomId,
		"isAvailable": false,
		"reservationId": reservation.Id,	
	}
	messageBody, _ := json.Marshal(message)
	err = rc.RabbitMQChannel.Publish(
		"", // exchange
		"reservation", // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body: messageBody,
		},
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": reservation})
}

func(rs *ReservationController) GetReservation(ctx *gin.Context)  {
	var id string = ctx.Param("id")
	reservation, err := rs.ReservationService.GetReservation(id)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, reservation)
}

func (rc *ReservationController) GetAllReservations(ctx *gin.Context)  {
	users, err := rc.ReservationService.GetAllReservations()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, users)
}
func (rc *ReservationController) Update(ctx *gin.Context)  {
	var reservation models.Reservation
	id := ctx.Param("id")
	if err := ctx.ShouldBindJSON(&reservation); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err := rc.ReservationService.Update(id,&reservation)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}
func (rc *ReservationController) Delete(ctx *gin.Context) {
	// Ekstrakcija ID-a iz URL parametra
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "ID is required"})
		return
	}

	// Poziv servisa za brisanje po ID-u
	err := rc.ReservationService.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	// Uspešan odgovor
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}
func (rc *ReservationController) GetGuestReservations(ctx *gin.Context) {
	// Ekstrakcija ID-a gosta iz URL parametra
	guestId,err := strconv.Atoi(ctx.Param("guest_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "guest_id is required!"})
		return
	}

	// Poziv servisa za dobavljanje rezervacija po ID-u gosta
	reservations, err := rc.ReservationService.GetGuestReservations(guestId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	// Uspešan odgovor
	ctx.JSON(http.StatusOK, reservations)
	
	
}


func(rc *ReservationController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/reservations", rc.Create)
	router.GET("/reservations/:id", rc.GetReservation)
	router.GET("/reservations", rc.GetAllReservations)
	router.PUT("/reservations/:id", rc.Update)
	router.DELETE("/reservations/:id", rc.Delete)
	router.GET("/guests/:guest_id/reservations", rc.GetGuestReservations)
}