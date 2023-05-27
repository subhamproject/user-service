package usrmgr

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/subhamproject/user-service/logs"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func CreateUserHandler(c *gin.Context) {

	tracer := otel.Tracer("CreateUserHandlerTrace")
	_, span := tracer.Start(c.Request.Context(), "CreateUserHandler")

	defer span.End()

	logs.DebugTrace(c.Request.Context(), span, "received request to create new user")
	var user User
	err := c.BindJSON(&user)
	if err != nil {
		logs.ErrorTrace(c.Request.Context(), span, fmt.Sprintf("unable parse create user request, error - %v", err))
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
		logs.ErrorTrace(c.Request.Context(), span, fmt.Sprintf("failed create user request, error - %v", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, usrId)
}

func GetAllUsersHandler(c *gin.Context) {
	tracer := otel.Tracer("GetAllUsersHandlerTrace")
	_, span := tracer.Start(c.Request.Context(), "GetAllUsersHandler")
	defer span.End()

	logs.DebugTrace(c.Request.Context(), span, "received request to get all users")
	users, err := GetAllUsers(c.Request.Context())
	if err != nil {
		logs.ErrorTrace(c.Request.Context(), span, fmt.Sprintf("failed to get users from db, error - %v", err))
		logs.Error(fmt.Sprintf("failed to get users from db, error - %v", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable retrieve users, error: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func GetUserHandler(c *gin.Context) {
	tracer := otel.Tracer("GetUserHandlerTrace")
	_, span := tracer.Start(c.Request.Context(), "GetUserHandler")
	defer span.End()
	id := c.Query("id")
	fmt.Printf("received request to get user by id %s\n", id)
	user, err := GetUserByID(c.Request.Context(), id)
	if err != nil {
		fmt.Printf("unable to get user by id %s , error - %v\n", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable retrieve user"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func GetUserOrderHandler(c *gin.Context) {
	tracer := otel.Tracer("GetUserOrderHandlerTrace")
	_, span := tracer.Start(c.Request.Context(), "GetUserOrderHandler")
	defer span.End()

	id := c.Query("id")
	fmt.Printf("received request to get user %s, orders \n", id)
	userOrder, err := GetUserOrder(c.Request.Context(), id)
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
