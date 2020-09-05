package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MonthEntry struct {
	ID    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	User  primitive.ObjectID `json:"user,omitempty" bson:"user,omitempty"`
	Year  int32              `json:"year,omitempty" bson:"year,omitempty"`
	Month int32              `json:"month,omitempty" bson:"month,omitempty"`
	Days  []DayEntry         `json:"days,omitempty" bson:"days,omitempty"`
}

type DayEntry struct {
	Week        int32              `json:"week,omitempty" bson:"week,omitempty"`
	DateTime    primitive.DateTime `json:"datetime,omitempty" bson:"datetime,omitempty"`
	Description string             `jsons:"description,omitempty" bson:"description,omitempty"`
	Priority    int32              `json:"priority,omitempty" bson:"priority,omitempty"`
	IsPermanent bool               `json:"ispermanent,omitempty" bson:"ispermanent,omitempty"`
}

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Email    string             `json:"email,omitempty" bson:"email,omitempty"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
}

type Config struct {
	Username     string `json:"dbUser"`
	Password     string `json:"dbPass"`
	ClusterName  string `json:"cluster"`
	CloudService string `json:"cloudService"`
}

var client *mongo.Client

//User
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

//Calendar Entry
func CreateMonthEntryEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")

	var monthEntry MonthEntry

	if err := json.NewDecoder(request.Body).Decode(&monthEntry); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	userCollection := client.Database("db_1").Collection("entries")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	result, err := userCollection.InsertOne(ctx, monthEntry)

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	json.NewEncoder(response).Encode(result)
}

func CreateDayEntryEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	param := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(param["id"])
	var dayEntry DayEntry

	if err := json.NewDecoder(request.Body).Decode(&dayEntry); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	if dayEntry.DateTime == 0 {
		dayEntry.DateTime = primitive.NewDateTimeFromTime(time.Now())
	}

	entryCollection := client.Database("db_1").Collection("entries")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	monthEntry := bson.M{"user": id}
	dayPush := bson.M{"$push": bson.M{"days": dayEntry}}
	result, err := entryCollection.UpdateOne(ctx, monthEntry, dayPush)

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	json.NewEncoder(response).Encode(result)
}

func GetEntryEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	param := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(param["id"])
	var entry MonthEntry
	collection := client.Database("db_1").Collection("entries")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, MonthEntry{User: id}).Decode(&entry)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
}
func EditEntryEndpoint(response http.ResponseWriter, request *http.Request) {

}
func DeleteEntryEndpoint(response http.ResponseWriter, request *http.Request) {

}

func LoadConfiguration() (Config, error) {
	var config Config
	configFile, err := os.Open("config.json")
	defer configFile.Close()
	if err != nil {
		return config, err
	}
	parse := json.NewDecoder(configFile)
	err = parse.Decode(&config)
	return config, err
}

func main() {
	var err error
	config, err := LoadConfiguration()
	if err != nil {
		log.Fatal(err)
		return
	}

	client, err = mongo.NewClient(options.Client().ApplyURI("mongodb+srv://" + config.Username + ":" + config.Password + "@" + config.ClusterName + "." + config.CloudService + ".mongodb.net/db?retryWrites=true&w=majority"))
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
	router.HandleFunc("/entry/month", CreateMonthEntryEndpoint).Methods("POST")
	router.HandleFunc("/entry/day/{id}", CreateDayEntryEndpoint).Methods("POST")
	router.HandleFunc("/entry/{id}", GetEntryEndpoint).Methods("GET")
	http.ListenAndServe(":12345", router)
}
