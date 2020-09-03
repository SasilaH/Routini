package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"

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
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Email    string             `json:"email,omitempty" bson:"email,omitempty"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
}

var client *mongo.Client

func UserSignupEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")

	var user User

	if err := json.NewDecoder(request.Body).Decode(&user); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	userCollection := client.Database("db_1").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	user.Password = string(hash)

	result, err := userCollection.InsertOne(ctx, user)

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	json.NewEncoder(response).Encode(result)
}

func UserLoginEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")

	var inputUser User

	if err := json.NewDecoder(request.Body).Decode(&inputUser); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	var returnUser User

	userCollection := client.Database("db_1").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	if err := userCollection.FindOne(ctx, bson.M{"email": inputUser.Email}).Decode(&returnUser); err != nil {
		log.Fatal(err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(returnUser.Password), []byte(inputUser.Password)); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	json.NewEncoder(response).Encode(returnUser)
}

func main() {

	var err error
	client, err = mongo.NewClient(options.Client().ApplyURI("mongodb+srv://userS:Teemoteemo123@cluster0.ge7wb.azure.mongodb.net/db?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// databases, err := client.ListDatabaseNames(ctx, bson.M{})
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(databases)

	router := mux.NewRouter()
	router.HandleFunc("/user", UserSignupEndpoint).Methods("POST")
	router.HandleFunc("/user", UserLoginEndpoint).Methods("GET")
	http.ListenAndServe(":12345", router)
}
