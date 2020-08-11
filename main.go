package main

import (
	dbf "app/db-go"
	rh "app/handlers-go"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	/// setup router
	r := mux.NewRouter()
	method := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origin := handlers.AllowedOrigins([]string{"*"})
	headers := handlers.AllowedHeaders([]string{"X-Request-With", "Content-Type", "Authorization"})
	db, err := dbf.RunDB()
	if err != nil {
		panic(err)
	}
	env := &dbf.Env{DB: db}
	r.Handle("/register", rh.RegisterHandler(env)).Methods("POST")
	r.Handle("/login", rh.LoginHandler(env)).Methods("POST")
	r.Handle("/ws", rh.WsHandler(env)).Methods("GET")
	//runing and listening to server on port 8000
	log.Fatal(http.ListenAndServe(":8000", handlers.CORS(headers, method, origin)(r)))
}
