// main.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sample-web-api/helpers"
	"sample-web-api/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var Objects []models.Object

var ctx = helpers.ConnectDB()

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func returnAllObjects(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllObjects")

	raw, err := ctx.Find(context.TODO(), bson.M{})

	if err != nil {
		helpers.GetError(err, w)
		return
	}

	defer raw.Close(context.TODO())

	for raw.Next(context.TODO()) {

		var obj models.Object

		err := raw.Decode(&obj)
		if err != nil {
			log.Fatal(err)
		}

		Objects = append(Objects, obj)
	}

	json.NewEncoder(w).Encode(Objects)
}

func returnSingleObject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var obj models.Object

	id, _ := primitive.ObjectIDFromHex(vars["id"])

	filter := bson.M{"_id": id}
	err := ctx.FindOne(context.TODO(), filter).Decode(&obj)

	if err != nil {
		helpers.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(obj)
}

func createNewObject(w http.ResponseWriter, r *http.Request) {

	var obj models.Object

	// we decode our body request params
	_ = json.NewDecoder(r.Body).Decode(&obj)

	// insert our book model.
	result, err := ctx.InsertOne(context.TODO(), obj)

	if err != nil {
		helpers.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func updateObject(w http.ResponseWriter, r *http.Request) {

	var vars = mux.Vars(r)

	//Get id from parameters
	id, _ := primitive.ObjectIDFromHex(vars["id"])

	var obj models.Object

	// Create filter
	filter := bson.M{"_id": id}

	// Read update model from body request
	_ = json.NewDecoder(r.Body).Decode(&obj)

	// prepare update model.
	update := bson.D{
		{"$set", bson.D{
			{"title", obj.Title},
			{"body", obj.Body},
			{"desc", obj.Desc},
		}},
	}

	err := ctx.FindOneAndUpdate(context.TODO(), filter, update).Decode(&obj)

	if err != nil {
		helpers.GetError(err, w)
		return
	}

	obj.Id = id

	json.NewEncoder(w).Encode(obj)
}

func deleteObject(w http.ResponseWriter, r *http.Request) {

	var vars = mux.Vars(r)

	// string to primitve.ObjectID
	id, err := primitive.ObjectIDFromHex(vars["id"])

	// prepare filter.
	filter := bson.M{"_id": id}

	result, err := ctx.DeleteOne(context.TODO(), filter)

	if err != nil {
		helpers.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(result)

}

func handleRequests() {
	config := helpers.GetConfiguration()

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/objects", returnAllObjects)
	myRouter.HandleFunc("/objects/add", createNewObject).Methods("POST")
	myRouter.HandleFunc("/objects/update", updateObject).Methods("PATCH")
	myRouter.HandleFunc("/objects/remove", deleteObject).Methods("DELETE")
	myRouter.HandleFunc("/objects/{id}", returnSingleObject)
	log.Fatal(http.ListenAndServe(config.Port, myRouter))
}

func main() {
	handleRequests()
}
