package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Entry struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	User        primitive.ObjectID `json:"user,omitempty" bson:"user,omitempty"`
	Date        string             `json:"date,omitempty" bson:"date,omitempty"`
	Time        string             `json:"time,omitempty" bson:"time,omitempty"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	Weight      int32              `json:"weight,omitempty" bson:"weight,omitempty"`
	Tags        []string           `json:"tags,omitempty" bson:"tags,omitempty"`
}

type User struct {
	ID   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name,omitempty" bson:"name,omitempty"`
}

var client *mongo.Client

func GenerateUserEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var user User
	if err := json.NewDecoder(request.Body).Decode(&user); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	userCollection := client.Database("db_1").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(result)
}

func main() {
	fmt.Println("111111")

	client, _ = mongo.NewClient(options.Client().ApplyURI("mongodb+srv://userS:Teemoteemo123@cluster0.ge7wb.azure.mongodb.net/db?retryWrites=true&w=majority"))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	fmt.Println("222222")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	fmt.Println("3333333")

	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		panic(err)
	}
	fmt.Println(databases)

	router := mux.NewRouter()
	router.HandleFunc("/user", GenerateUserEndpoint).Methods("POST")
	http.ListenAndServe(":12345", router)
}
