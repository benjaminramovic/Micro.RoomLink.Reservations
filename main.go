package main

import (
	"context"
	"example/micro-roomlink-reservations/controllers"
	"example/micro-roomlink-reservations/services"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

/*
	@Author: DevProblems(Sarang Kumar)
	@YTChannel: https://www.youtube.com/channel/UCVno4tMHEXietE3aUTodaZQ
*/
var (
	server      *gin.Engine
	rs          services.ReservationService
	rc          *controllers.ReservationController
	ctx         context.Context
	reservations       *mongo.Collection
	mongoclient *mongo.Client
	err         error
)

func init() {
	ctx = context.TODO()

	mongoconn := options.Client().ApplyURI("mongodb://localhost:27017")
	mongoclient, err = mongo.Connect(ctx, mongoconn)
	if err != nil {
		log.Fatal("error while connecting with mongo", err)
	}
	err = mongoclient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("error while trying to ping mongo", err)
	}

	fmt.Println("mongo connection established")

	reservations = mongoclient.Database("reservationsdb").Collection("reservations")
	rs = services.NewReservationService(reservations, ctx)
	rc = controllers.NewReservationController(rs)
	server = gin.Default()
}

func main() {
	defer mongoclient.Disconnect(ctx)

	basepath := server.Group("/v1")
	rc.RegisterRoutes(basepath)

	log.Fatal(server.Run(":9090"))

}

