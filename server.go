package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type User struct {
	Name string `json:"name"`
}

func NewUser(name string) User {
	return User{Name: name}
}

var admins []User

func init() {
	admins = append(admins, NewUser("alice"), NewUser("bob"))
}

func main() {
	http.HandleFunc("/admins", func(w http.ResponseWriter, r *http.Request) {
		js, err := json.Marshal(admins)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	})
	log.Println("listening to port *:8080. press ctrl + c to cancel")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
