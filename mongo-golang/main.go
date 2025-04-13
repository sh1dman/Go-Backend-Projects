package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/shadman/mongo-golang/controllers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	r := httprouter.New()
	uc := controllers.NewUserController(getSession())

	r.GET("/user/:id", uc.GetUser)
	r.POST("/user", uc.CreateUser)       // Fixed route
	r.DELETE("/user/:id", uc.DeleteUser) // Fixed route

	http.ListenAndServe("localhost:8080", r)
}

func getSession() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		"mongodb://admin:password123@localhost:27017/?authSource=admin",
	))
	if err != nil {
		panic(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to MongoDB!")
	return client
}
