package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/Kibuns/TwootService/DAL"
	"github.com/Kibuns/TwootService/Models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	fmt.Println("Twoot service started! yipeee")
	handleRequests()
}

//controllerssss

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage of the Twoot Service!")
	fmt.Println("Endpoint Hit: homePage")
}

func returnAll(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAll")
	json.NewEncoder(w).Encode(getAllTwoots())
}

func returnTwoot(w http.ResponseWriter, r *http.Request) {
	var idParam string = mux.Vars(r)["id"]
	json.NewEncoder(w).Encode(DAL.ReadSingleTwoot(idParam))
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	var idParam string = mux.Vars(r)["user"]
	DAL.DeleteAllOfUser(idParam)
	fmt.Fprintf(w, "deleted everything from user: " + idParam)
}

func storeTwoot(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	fmt.Println("Storing Twoot")
	// parse the request body into a Twoot struct
	var twoot Models.Twoot
	err := json.NewDecoder(body).Decode(&twoot)
	fmt.Println(twoot.Content)
	if err != nil {
		http.Error(w, "Could not decode body into twoot", http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	// check if the twoot content contains a hashtag
	if strings.Contains(twoot.Content, "#") {
		// use a regular expression to check if the hashtag is followed by any letter
		match, _ := regexp.MatchString(`#\p{L}+`, twoot.Content)
		if match {
			fmt.Println("hashtag found in content")
			send("twoot", &twoot) //message to search service
		} else {
			fmt.Println("no hashtag found")
		}
	}

	// insert the twoot into the database
	DAL.InsertTwoot(twoot, w)

}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.Use(CORS)

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/all", returnAll)
	myRouter.HandleFunc("/get/{id}", returnTwoot)
	myRouter.HandleFunc("/create", storeTwoot)
	myRouter.HandleFunc("/delete/{user}", deleteUser)

	log.Fatal(http.ListenAndServe(":10000", myRouter))
}



//service functions

func getAllTwoots() (values []primitive.M) {
	return DAL.ReadAllTwoots()
}

// other
// CORS Middleware
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Set headers
		w.Header().Set("Access-Control-Allow-Headers:", "*")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		fmt.Println("ok")

		// Next
		next.ServeHTTP(w, r)
		//return
	})

}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
