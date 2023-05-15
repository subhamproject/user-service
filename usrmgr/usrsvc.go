package usrmgr

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type User struct {
	ID    string      `json:"id"`
	Name  string      `json:"name"`
	Order interface{} `json:"order,omitempty"`
}

func GetUserByID(ctx context.Context, id string) (User, error) {

	SendLogs(fmt.Sprintf("received request to get user by Id %s", id))

	var user User
	filter := bson.D{{"id", id}}
	err := userCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func GetAllUsers() ([]User, error) {

	SendLogs("received request to get all users")
	var users []User
	cursor, err := userCollection.Find(context.Background(), bson.D{{}})
	if err != nil {
		return users, err
	}

	if err = cursor.All(context.TODO(), &users); err != nil {
		return users, err
	}

	return users, nil
}

func CreateUser(ctx context.Context, usr User) (string, error) {

	// tr := otel.Tracer("CreateUser")
	// _, span := tr.Start(ctx, "CreateUser-service")
	// span.SetAttributes(attribute.Key("CreateUser").String(usr.Name))
	// defer span.End()
	// SendLogs(fmt.Sprintf("received request to create new user: %s", usr.Name))

	tracer := otel.Tracer("CreateUserServiceTrace")
	_, span := tracer.Start(ctx, "CreateUserService")

	defer span.End()

	// get the current span by the request context
	currentSpan := trace.SpanFromContext(ctx)
	currentSpan.AddEvent("CreateUserService-Event")
	currentSpan.SetAttributes(attribute.String("UserName", usr.Name))

	id := genUserId()
	usr.ID = id

	result, err := userCollection.InsertOne(ctx, usr)
	if err != nil {
		fmt.Printf("failed to insert user: %v\n", err)
		return "", err
	}
	fmt.Println("user inserted with InsertedID: ", result.InsertedID)

	CreateUserOrder(ctx, usr.ID)

	return id, nil
}

func GetUserOrder(ctx context.Context, id string) (User, error) {
	host := GetEnvParam("ORDER_SVC_HOST", "localhost")
	port := GetEnvParam("ORDER_SVC_PORT", "8081")

	orderSvcUrl := fmt.Sprintf("http://%s:%s/order?userId=%s", host, port, id)
	resp, err := http.Get(orderSvcUrl)
	if err != nil {
		fmt.Println("error while loading user orders ", err)
		return User{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return User{}, fmt.Errorf("error while reading user order response- %v", err)
	}
	var order interface{}
	err = json.Unmarshal(body, &order)
	if err != nil {
		fmt.Println("error while parsing user orders ", err)
		return User{}, fmt.Errorf("error while parsing user orders- %v", err)
	}
	usr, err := GetUserByID(ctx, id)
	if err != nil {
		fmt.Println("error while finding user by Id", err)
		return User{}, err
	}

	usr.Order = order
	return usr, nil
}

func CreateUserOrder(ctx context.Context, userId string) error {

	// tr := otel.Tracer("CreateUserOrder")
	// _, span := tr.Start(ctx, "CreateUserOrder-service")
	// span.SetAttributes(attribute.Key("CreateUserOrder-UserId").String(userId))
	// defer span.End()

	// get the current span by the request context
	currentSpan := trace.SpanFromContext(ctx)
	currentSpan.AddEvent("CreateUserOrder-Event")
	currentSpan.SetAttributes(attribute.String("UserId", userId))

	tracer := otel.Tracer("CreateUserOrderTrace")
	_, span := tracer.Start(ctx, "CreateUserOrder")

	defer span.End()

	fmt.Printf("invoke order-service to create order for user %v\n", userId)

	host := GetEnvParam("ORDER_SVC_HOST", "localhost")
	port := GetEnvParam("ORDER_SVC_PORT", "8081")

	orderSvcUrl := fmt.Sprintf("http://%s:%s/order?userId=%s", host, port, userId)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, orderSvcUrl, bytes.NewBuffer(nil))
	if err != nil {
		fmt.Println("failed to creare request for user orders ", err)
		return err
	}
	httpClient := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("failed to creare user orders ", err)
		return err
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println("failed to creare user orders received response ", resp)
		return err
	}
	return nil
}

// GetEnvParam : return string environmental param if exists, otherwise return default
func GetEnvParam(param string, dflt string) string {
	if v, exists := os.LookupEnv(param); exists {
		return v
	}
	return dflt
}

func genUserId() string {
	// Create a big.Int with the maximum value for the desired range
	max := big.NewInt(10000)

	// Generate a random big.Int
	// The first argument is a reader that returns random numbers
	// The second argument is the maximum value (not inclusive)
	randInt, err := rand.Int(rand.Reader, max)

	if err != nil {
		fmt.Println("Error generating random number:", err)
		return "100"
	}
	return randInt.String()
}
