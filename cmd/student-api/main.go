package main

import (
	"context"
	"fmt"
	"github.com/hassanjawwad12/student-api/internal/config"
	"github.com/hassanjawwad12/student-api/internal/http/handlers/student"
	"github.com/hassanjawwad12/student-api/internal/storage/sqlite"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fmt.Println("Welcome")
	//load config
	cfg := config.MustLoad()

	//db setup
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("storage initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	//setup router
	//servemux stores a mapping between the predefined URL paths for your application and the corresponding handlers.
	//Usually you have one servemux for your application containing all your routes.
	router := http.NewServeMux()
	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))

	//setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}
	slog.Info("Server started", slog.String("address:", cfg.Addr))

	//Implementation of graceful shutdown
	//Completes the  on-going request but does not entertain the incoming ones

	//size set to 1 because buffer channel
	done := make(chan os.Signal, 1)

	//notify us on every interrupt signal
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to start server")
		}
	}()

	//unless a signal enters done , our code won't move forward
	<-done

	slog.Info("Shutting down the server")

	//added a 5 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//we will throw error after 5 seconds
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("server shutdown successfully")

}
