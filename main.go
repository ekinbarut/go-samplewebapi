// main.go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// Object - Our struct for all objects
type Object struct {
	Id    string `json:"Id"`
	Title string `json:"Title"`
	Desc  string `json:"desc"`
	Body  string `json:"body"`
}

var Objects []Object

func readJsonFile() []Object {
	raw, err := ioutil.ReadFile("./objects.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	json.Unmarshal(raw, &Objects)
	return Objects
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func returnAllObjects(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllObjects")
	json.NewEncoder(w).Encode(Objects)
}

func returnSingleObject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	for _, object := range Objects {
		if object.Id == key {
			json.NewEncoder(w).Encode(object)
		}
	}
}

func createNewObject(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var object Object
	json.Unmarshal(reqBody, &object)

	if !exists(object, Objects) {

		Objects = append(Objects, object)

		json.NewEncoder(w).Encode(object)
		file, _ := json.MarshalIndent(Objects, "", "")
		_ = ioutil.WriteFile("objects.json", file, 0644)
	} else {
		json.NewEncoder(w).Encode("Id already exists.")
	}
}

func updateObject(w http.ResponseWriter, r *http.Request) {

	i := 0

	reqBody, _ := ioutil.ReadAll(r.Body)
	var object Object
	json.Unmarshal(reqBody, &object)

	if exists(object, Objects) {

		//find index of the existing value
		for index, obj := range Objects {
			if obj.Id == object.Id {
				i = index
			}
		}

		Objects[i].Body = object.Body
		Objects[i].Desc = object.Desc
		Objects[i].Title = object.Title

		json.NewEncoder(w).Encode(object)
		file, _ := json.MarshalIndent(Objects, "", "")
		_ = ioutil.WriteFile("objects.json", file, 0644)
	} else {
		json.NewEncoder(w).Encode("Id does not exist, can not update.")
	}
}

func deleteObject(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var object Object
	json.Unmarshal(reqBody, &object)

	if exists(object, Objects) {

		for index, obj := range Objects {
			if obj.Id == object.Id {
				Objects = append(Objects[:index], Objects[index+1:]...)
			}
		}

		json.NewEncoder(w).Encode(object)
		file, _ := json.MarshalIndent(Objects, "", "")
		_ = ioutil.WriteFile("objects.json", file, 0644)
	} else {
		json.NewEncoder(w).Encode("Id does not exist, can not")
	}

}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/objects", returnAllObjects)
	myRouter.HandleFunc("/objects/add", createNewObject).Methods("POST")
	myRouter.HandleFunc("/objects/update", updateObject).Methods("PATCH")
	myRouter.HandleFunc("/objects/remove", deleteObject).Methods("DELETE")
	myRouter.HandleFunc("/objects/{id}", returnSingleObject)
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	readJsonFile()
	handleRequests()
}
