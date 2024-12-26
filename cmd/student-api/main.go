package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hassanjawwad12/student-api/internal/config"
)

func main() {
	fmt.Println("Welcome")
	//load config
	cfg := config.MustLoad()

	//db setup
	//setup router
	//servemux stores a mapping between the predefined URL paths for your application and the corresponding handlers.
	//Usually you have one servemux for your application containing all your routes.
	router := http.NewServeMux()
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		// Convert the string to a byte slice before writing
		w.Write([]byte("Welcome to student API"))
	})
	//setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}
	fmt.Println("Server started")

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Failed to start server")
	}

}
