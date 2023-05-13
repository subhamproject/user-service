package usrmgr

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func CreateUserHandler(c *gin.Context) {

	tracer := otel.Tracer("CreateUserHandlerTrace")
	_, span := tracer.Start(c.Request.Context(), "CreateUserHandler")

	defer span.End()

	// log.Printf("In CreateUserHandler span, after calling a child function. When this function ends, parentSpan will complete.")

	fmt.Println("received request to create new user")
	var user User
	err := c.BindJSON(&user)
	if err != nil {
		fmt.Printf("unable parse create user request, error - %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.SetAttributes(attribute.String("UserName", user.Name))

	// get the current span by the request context
	currentSpan := trace.SpanFromContext(c.Request.Context())
	currentSpan.AddEvent("CreateUserHandler-Event")
	currentSpan.SetAttributes(attribute.String("UserName", user.Name))

	usrId, err := CreateUser(c.Request.Context(), user)
	if err != nil {
		fmt.Printf("failed create user request, error - %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, usrId)
}

func GetAllUsersHandler(c *gin.Context) {
	fmt.Println("received request to get all users")
	users, err := GetAllUsers()
	if err != nil {
		fmt.Printf("failed to get users from db, error - %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable retrieve users, error: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func GetUserHandler(c *gin.Context) {
	id := c.Query("id")
	fmt.Printf("received request to get user by id %s\n", id)
	user, err := GetUserByID(c, id)
	if err != nil {
		fmt.Printf("unable to get user by id %s , error - %v\n", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable retrieve user"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func GetUserOrderHandler(c *gin.Context) {
	id := c.Query("id")
	fmt.Printf("received request to get user %s, orders \n", id)
	userOrder, err := GetUserOrder(c, id)
	if err != nil {
		fmt.Printf("unable to get user %s, orders. error - %v \n", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("unable retrieve user's order data, %v", err)})
		return
	}
	c.JSON(http.StatusOK, userOrder)
}

func GetServiceHealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, "I'm Healthly")
}
