package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type user struct {
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	Age       int       `json:"age"`
	ID        uuid.UUID `json:"id"`
}

type userInput struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Age       int    `json:"age"`
}

var users map[uuid.UUID]user

func withHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func handleGetAllUsers(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("content-type", "application/json")
	err := json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var userInput userInput

	err := json.NewDecoder(r.Body).Decode(&userInput)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	log.Println(userInput)

	user := user{
		Firstname: userInput.Firstname,
		Lastname:  userInput.Lastname,
		Age:       userInput.Age,
		ID:        uuid.New(),
	}

	log.Println(user)

	users[user.ID] = user

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func handleGetUser(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	user, ok := users[id]
	if !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	err := json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func handleUpdateUser(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	var userInput userInput

	err := json.NewDecoder(r.Body).Decode(&userInput)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
	user := user{
		Firstname: userInput.Firstname,
		Lastname:  userInput.Lastname,
		Age:       userInput.Age,
		ID:        id,
	}
	users[user.ID] = user
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func handleDeleteUser(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	_, ok := users[id]
	if !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	delete(users, id)
	w.WriteHeader(http.StatusNoContent)
}

func parseUUIDFromRequest(r *http.Request) (uuid.UUID, error) {
	vars := mux.Vars(r)
	return uuid.Parse(vars["id"])
}

//https://users

func homeHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	withHeaders(w)

	switch r.Method {
	case http.MethodGet:
		handleGetAllUsers(w, r)
		return
	case http.MethodPost:
		handleCreateUser(w, r)
		return
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	withHeaders(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	parsedId, err := parseUUIDFromRequest(r)
	if err != nil {
		http.Error(w, "Неверный UUID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		handleGetUser(w, r, parsedId)
		return
	case http.MethodPut:
		handleUpdateUser(w, r, parsedId)
		return
	case http.MethodDelete:
		handleDeleteUser(w, r, parsedId)
		return
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func main() {
	r := mux.NewRouter()
	users = make(map[uuid.UUID]user)
	r.HandleFunc("/users", homeHandler).Methods("GET", "OPTIONS", "POST")
	r.HandleFunc("/users/{id}", usersHandler).Methods("GET", "OPTIONS", "PUT", "DELETE")
	http.ListenAndServe(":8080", r)
}
