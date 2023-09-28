package helpers

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectSessionDBMongo() *mongo.Client {
	user := "surnet"
	password := "APIsurnet2008"
	direction := "192.168.20.150:27017"
	direction2 := "192.168.20.151:27017"
	timeoutMs := "5000"

	//stringLink := fmt.Sprintf("mongodb://%s:%s@%s,%s/?serverSelectionTimeoutMS=%s&authSource=admin&replicaSet=snt0", user, password, direction, direction2, timeoutMs)
	stringLink := fmt.Sprintf("mongodb://%s:%s@%s, %s/?serverSelectionTimeoutMS=%s&authSource=admin&replicaSet=snt0", user, password, direction, direction2, timeoutMs)

	// fmt.Println("StringLink: ", stringLink)

	clientOptions := options.Client().ApplyURI(stringLink)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	return client
}
