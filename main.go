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
	"github.com/subhamproject/user-service/otelsvc"
	"github.com/subhamproject/user-service/usrmgr"
	"github.com/subhamproject/user-service/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

var (
	client       *mongo.Client
	ctx          context.Context
	cFund        context.CancelFunc
	server       *http.Server
	otelShutdown func()
)

func main() {

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		fmt.Println("initializing connection mongo...")
		//init mogno db
		client, ctx, cFund, _ = usrmgr.InitMongoDB()

		//init kafka connection
		usrmgr.InitKafka()
	}()

	wg.Wait()

	log.Printf("initializing otel connection...")
	otelUrl := utils.GetEnvParam("OTEL_COLLECTOR_URL", "localhost:4317")
	otelEnable := utils.GetEnvBoolParam("OTEL_ENABLE", false)
	otelShutdown = otelsvc.InitTracerProvider(otelUrl, otelEnable)

	r := gin.Default()

	f := func(req *http.Request) bool { return req.URL.Path != "/health" }
	r.Use(otelgin.Middleware("user-service", otelgin.WithFilter(f)))

	r.GET("/health", usrmgr.GetServiceHealthHandler)
	r.POST("/user", usrmgr.CreateUserHandler)
	r.GET("/users", usrmgr.GetAllUsersHandler)
	r.GET("/user", usrmgr.GetUserHandler)
	r.GET("/user/order", usrmgr.GetUserOrderHandler)

	serverPort := utils.GetEnvParam("SERVICE_PORT", "8082")

	server = &http.Server{
		Addr:    ":" + serverPort,
		Handler: r,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	shutdownServer()

	log.Println("Server exiting")

}

func shutdownServer() {
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	//close mongo driver
	usrmgr.CloseMongoDB(client, ctx, cFund)
	//close kafka connection
	usrmgr.CloseKafka()

	otelShutdown()

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

}
