package usrmgr

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateUserHandler(c *gin.Context) {
	fmt.Println("received request to create new user")
	var user User
	err := c.BindJSON(&user)
	if err != nil {
		fmt.Printf("unable parse create user request, error - %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	usrId, err := CreateUser(user)
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
	user, err := GetUserByID(id)
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
	userOrder, err := GetUserOrder(id)
	if err != nil {
		fmt.Printf("unable to get user %s, orders. error - %v \n", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("unable retrieve user's order data, %v", err)})
		return
	}
	c.JSON(http.StatusOK, userOrder)
}
