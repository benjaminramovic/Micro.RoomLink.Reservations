package main

import (
	"context"
	"example/micro-roomlink-reservations/controllers"
	"example/micro-roomlink-reservations/services"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	//"github.com/rabbitmq/amqp091-go"
	"github.com/streadway/amqp"
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
	rabbitChannel *amqp.Channel
	rabbitConn *amqp.Connection
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

	//rabbitMQ coonection setup
	rabbitConn,err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println("Successfully connected to RabbitMQ")

	rabbitChannel, err := rabbitConn.Channel()
	if err != nil {
		fmt.Println(err)
		panic(err)
	
	}


	  // Provera da li je RabbitMQ veza otvorena
	  if rabbitConn.IsClosed() {
        fmt.Println("RabbitMQ connection is closed, reconnecting...")
        rabbitConn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
        if err != nil {
            log.Fatal("Failed to reconnect to RabbitMQ:", err)
        }
    }

    // Provera da li je kanal otvoren
    if rabbitChannel == nil {
        fmt.Println("Channel is nil or closed")
        return
    }

	q, err := rabbitChannel.QueueDeclare(
		"reservation", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		fmt.Println(err)
		panic(err)

	}
	fmt.Println("Queue created: ", q)


	rc = controllers.NewReservationController(rs,rabbitChannel)
	server = gin.Default()
}

func main() {
	
	/*err = ch.Publish(
		"",
		"reservation",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("Hello World"),
		},
	)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("Message published")*/

	

	defer func() {
		if err := rabbitChannel.Close(); err != nil {
			fmt.Println("Failed to close RabbitMQ channel:", err)
		}
		if err := rabbitConn.Close(); err != nil {
			fmt.Println("Failed to close RabbitMQ connection:", err)
		}
	}()


	defer mongoclient.Disconnect(ctx)

	basepath := server.Group("/api")
	rc.RegisterRoutes(basepath)

	log.Fatal(server.Run(":9090"))

}

