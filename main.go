package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type URLShortner struct {
	id        int
	url       string
	count     int
	createdAt time.Time
}

type CreateRequestBody struct {
	Url string `json:"url"`
}

type CreateRequestResponse struct {
	Id int `json:"id"`
}

var count = 0

var data = make(map[int]URLShortner)

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthCheck", HealthHandler)
	mux.HandleFunc("POST /create", CreateHandler)
	mux.HandleFunc("GET /redirect/{id}", RedirectHandler)

	// Start the web server on port 8080
	fmt.Println("Server listening on port 8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Printf("Server failed: %s\n", err)
	}
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	var body CreateRequestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	newShortener := ShortnerFactory(body.Url)
	data[newShortener.id] = *newShortener
	response := CreateRequestResponse{
		Id: newShortener.id,
	}
	count++
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201
	json.NewEncoder(w).Encode(response)
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {

	Id, _ := strconv.Atoi(r.PathValue("id"))

	urlShortnerInstance, exists := data[Id]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	redirectURL := urlShortnerInstance.url
	if !strings.HasPrefix(redirectURL, "http://") && !strings.HasPrefix(redirectURL, "https://") {
		redirectURL = "https://" + redirectURL
	}

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)

}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Webserver alive on port 8080")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Health Check"))
}

func ShortnerFactory(url string) *URLShortner {
	return &URLShortner{
		id:        count,
		url:       url,
		count:     count,
		createdAt: time.Now(),
	}
}
