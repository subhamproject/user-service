package usrmgr

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/subhamproject/user-service/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)

var userCollection *mongo.Collection

// This is a user defined method to close resources.
// This method closes mongoDB connection and cancel context.
func close(client *mongo.Client, ctx context.Context,
	cancel context.CancelFunc) {

	// CancelFunc to cancel to context
	defer cancel()

	// client provides a method to close
	// a mongoDB connection.
	defer func() {

		// client.Disconnect method also has deadline.
		// returns error if any,
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

// This is a user defined method that returns mongo.Client,
// context.Context, context.CancelFunc and error.
// mongo.Client will be used for further database operation.
// context.Context will be used set deadlines for process.
// context.CancelFunc will be used to cancel context and
// resource associated with it.
func connectLocal(uri, user, pass string) (*mongo.Client, context.Context,
	context.CancelFunc, error) {

	// ctx will be used to set deadline for process, here
	// deadline will of 30 seconds.
	ctx, cancel := context.WithTimeout(context.Background(),
		30*time.Second)

	credential := options.Credential{
		Username: user,
		Password: pass,
	}
	opts := options.Client().ApplyURI(uri).SetAuth(credential)

	//Add instrumentation to client options
	opts.Monitor = otelmongo.NewMonitor()
	// mongo.Connect return mongo.Client method
	client, err := mongo.Connect(ctx, opts)

	return client, ctx, cancel, err
}

func connect(uri string) (*mongo.Client, context.Context,
	context.CancelFunc, error) {

	// ctx will be used to set deadline for process, here
	// deadline will of 30 seconds.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	credential := options.Credential{
		AuthMechanism: "MONGODB-X509",
	}
	clientOpts := options.Client().ApplyURI(uri).SetAuth(credential)
	// Add instrumentation to client options
	clientOpts.Monitor = otelmongo.NewMonitor()

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		fmt.Println("connect failed!")
		log.Fatal(err)
		return nil, ctx, cancel, err
	}
	fmt.Println("connect successful!")

	return client, ctx, cancel, err
}

// This is a user defined method that accepts
// mongo.Client and context.Context
// This method used to ping the mongoDB, return error if any.
func ping(client *mongo.Client, ctx context.Context) error {

	// mongo.Client has Ping to ping mongoDB, deadline of
	// the Ping method will be determined by cxt
	// Ping method return error if any occurred, then
	// the error can be handled.
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		fmt.Println("ping failed!")
		log.Fatal(err)
		return err
	}
	fmt.Println("ping  successful!")
	return nil
}

func InitMongoDB() (*mongo.Client, context.Context,
	context.CancelFunc, error) {

	devMode := utils.GetEnvBoolParam("DEV_MODE", true)
	user := utils.GetEnvParam("MONGO_USERNAME", "root")
	pass := utils.GetEnvParam("MONGO_PASSWORD", "rootpassword")
	caFilePath := utils.GetEnvParam("MONGO_CA_CERT", "/home/om/go/src/github.com/subhamproject/devops-demo/certs/mongoCA.crt")
	certificateKeyFilePath := utils.GetEnvParam("MONGO_CLIENT_CERT_KEY", "/home/om/go/src/github.com/subhamproject/devops-demo/certs/mongo-client.pem")
	uri := utils.GetEnvParam("MONGO_URL", "mongodb://%s:%s@mongo1:27011,mongo2:27012,mongo3:27013/demo?replicaSet=rs0&tlsCAFile=%s&tlsCertificateKeyFile=%s")

	var client *mongo.Client
	var ctx context.Context
	var cFunc context.CancelFunc
	var err error

	if devMode {
		uri = "mongodb://localhost:27017"
		client, ctx, cFunc, err = connectLocal(uri, user, pass)
	} else {
		uri = fmt.Sprintf(uri, user, pass, caFilePath, certificateKeyFilePath)
		client, ctx, cFunc, err = connect(uri)
	}
	// Get Client, Context, CancelFunc and
	// err from connect method.
	if err != nil {
		panic(err)
	}

	// Ping mongoDB with Ping method
	ping(client, ctx)

	userCollection = client.Database("demo").Collection("users")

	return client, ctx, cFunc, err
}

func CloseMongoDB(client *mongo.Client, ctx context.Context, cancel context.CancelFunc) {

	// Release resource when the main
	// function is returned.
	close(client, ctx, cancel)

	fmt.Println("mongodb connected closed")
}
