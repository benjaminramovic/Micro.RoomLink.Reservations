package controllers

import (
	"example/micro-roomlink-reservations/models"
	"example/micro-roomlink-reservations/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReservationController struct {
	ReservationService services.ReservationService

}

func NewReservationController(reservationService services.ReservationService) *ReservationController {
	return &ReservationController{ReservationService: reservationService}
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
	ctx.JSON(http.StatusCreated, gin.H{"data": reservation})
}

func(rs *ReservationController) GetReservation(ctx *gin.Context)  {
	ctx.JSON(200, nil)
}

func (rc *ReservationController) GetAllReservations(ctx *gin.Context)  {
	users, err := rc.ReservationService.GetAllReservations()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, users)
}
func (rs *ReservationController) Update(ctx *gin.Context)  {
	ctx.JSON(200, nil)
}
func (rs *ReservationController) Delete(ctx *gin.Context)  {
	ctx.JSON(200, nil)
}

func(rc *ReservationController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/reservations", rc.Create)
	router.GET("/reservations/:id", rc.GetReservation)
	router.GET("/reservations", rc.GetAllReservations)
	router.PUT("/reservations/:id", rc.Update)
	router.DELETE("/reservations/:id", rc.Delete)
}