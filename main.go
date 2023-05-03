package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/subhamproject/user-service/usrmgr"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {

	var client *mongo.Client
	var ctx context.Context
	var cFund context.CancelFunc

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		fmt.Println("initializing connection mongo...")
		//init mogno db
		client, ctx, cFund, _ = usrmgr.InitMongoDB()
	}()

	wg.Wait()

	r := gin.Default()

	r.POST("/user", usrmgr.CreateUserHandler)
	r.GET("/users", usrmgr.GetAllUsersHandler)
	r.GET("/user", usrmgr.GetUserHandler)
	r.GET("/user/order", usrmgr.GetUserOrderHandler)

	serverPort := usrmgr.GetEnvParam("SERVICE_PORT", "8082")

	srv := &http.Server{
		Addr:    ":" + serverPort,
		Handler: r,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	//clode mongo driver
	usrmgr.CloseMongoDB(client, ctx, cFund)

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")

}
