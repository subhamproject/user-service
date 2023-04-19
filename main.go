package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
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
	fmt.Println(fmt.Sprintf("received request to get user by id %s", id))
	user, err := GetUserByID(id)
	if err != nil {
		fmt.Println(fmt.Sprintf("unable to get user by id %s , error - %v", id, err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable retrieve user"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func GetUserOrderHandler(c *gin.Context) {
	id := c.Query("id")
	fmt.Println(fmt.Sprintf("received request to get user %s, orders ", id))
	userOrder, err := GetUserOrder(id)
	if err != nil {
		fmt.Println(fmt.Sprintf("unable to get user %s, orders. error - %v", id, err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("unable retrieve user's order data, %v", err)})
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
	host := GetEnvParam("ORDER_SVC_HOST", "localhost")
	port := GetEnvParam("ORDER_SVC_PORT", "9092")
	orderSvcUrl := fmt.Sprintf("http://%s:%s/order?id=111", host, port)
	resp, err := http.Get(orderSvcUrl)
	if err != nil {
		fmt.Println("error while loading user orders ", err)
		return User{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	var order interface{}
	err = json.Unmarshal(body, &order)
	if err != nil {
		fmt.Println("error while parsing user orders ", err)
		return User{}, err
	}
	return User{
		ID:    id,
		Name:  "Doe",
		Order: order,
	}, nil
}

// GetEnvParam : return string environmental param if exists, otherwise return default
func GetEnvParam(param string, dflt string) string {
	if v, exists := os.LookupEnv(param); exists {
		return v
	}
	return dflt
}
