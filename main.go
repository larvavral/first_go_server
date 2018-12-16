package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"

	common "./common"
)

var server = "localhost"
var port = 1433
var user = "sa"
var password = "Hoangpq14"
var database = "simple_server"

var router = mux.NewRouter()

func main() {
	// Create connection string
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s",
		server, user, password, port, database)
	common.InitDB(connString)

	router.HandleFunc("/", common.LoginPageHandler)
	router.HandleFunc("/login", common.LoginHandler).Methods("POST")

	router.HandleFunc("/signup", common.SignupPageHandler).Methods("GET")
	router.HandleFunc("/signup", common.SignupHandler).Methods("POST")

	router.HandleFunc("/search", common.SearchPageHandler).Methods("GET")
	router.HandleFunc("/search", common.SearchHandler).Methods("POST")

	router.HandleFunc("/logout", common.LogoutHandler).Methods("GET")

	http.Handle("/", router)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.ListenAndServe(":8080", nil)
}