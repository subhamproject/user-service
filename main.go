package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type User struct {
	ID    string      `json:"id"`
	Name  string      `json:"name"`
	Order interface{} `json:"order,omitempty"`
}

func main() {
	r := gin.Default()
	r.GET("/user", GetUserHandler)
	r.GET("/user/order", GetUserOrderHandler)
	err := r.Run(":9091")
	if err != nil {
		log.Fatalf("impossible to start server: %s", err)
	}
}

func GetUserHandler(c *gin.Context) {
	id := c.Query("id")
	user, err := GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable retrieve user"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func GetUserOrderHandler(c *gin.Context) {
	id := c.Query("id")
	userOrder, err := GetUserOrder(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable retrieve user's order data"})
		return
	}
	c.JSON(http.StatusOK, userOrder)
}

func GetUserByID(id string) (User, error) {
	// TODO : lookup in db
	return User{
		ID:   id,
		Name: "Doe",
	}, nil
}

func GetUserOrder(id string) (User, error) {
	return User{
		ID:   id,
		Name: "Doe",
	}, nil
}
