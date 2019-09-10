package database

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var hostname string
var port string

// DBTimeout is the maximum response time from DB
const DBTimeout = 500

func setMongoParameters() {
	if os.Getenv("MONGO_HOSTNAME") != "" {
		hostname = os.Getenv("MONGO_HOSTNAME")
	} else {
		fmt.Print("USING LOCAL DATABASE")
		hostname = "localhost"
	}
	if os.Getenv("MONGO_PORT") != "" {
		port = os.Getenv("MONGO_PORT")
	} else {
		fmt.Print("USING DEFAULT DATABASE")
		port = "27017"
	}
	fmt.Print("DB: {name: mongo, hostname:" + hostname + ", port:" + port + "}")
}

func getClient(c *gin.Context) (*mongo.Client, error) {

	client, err := mongo.Connect(c, options.Client().ApplyURI("mongodb://"+hostname+":27017"))
	if err != nil {
		return nil, errors.New("Failed to generate Mongo Client")
	}

	fmt.Print("pinging database with this FQDN: " + hostname)

	// Short timeout to test mongo connection
	shortCtx, cancelFunc := context.WithTimeout(c, DBTimeout*time.Millisecond)
	defer cancelFunc()
	err = client.Ping(shortCtx, readpref.Primary())
	if err != nil {
		return nil, errors.New("Unable to reach database within " + strconv.Itoa(DBTimeout) + "ms")
	}
	fmt.Print("Acces granted !")
	return client, nil
}

func getDatabase(c *mongo.Client) *mongo.Database {
	name := "GoSmartSearchDatabase"
	database := c.Database(name)
	return database
}
