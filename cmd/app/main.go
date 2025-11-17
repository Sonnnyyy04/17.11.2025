package main

import (
	"context"
	"fmt"
	"links_project/internal/handler"
	"links_project/internal/service"
	"links_project/internal/storage"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	s := storage.NewStorage("storage.json")
	if err := s.Load(); err != nil{
		log.Fatalf("failed to load storage: %v", err)
	}
	srv := service.NewService(s)
	h := handler.NewHandler(srv)

	mux := http.NewServeMux()
	mux.HandleFunc("/links", h.HandleLinks)
	mux.HandleFunc("/report", h.HandleReport)

	server := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	go func(){
		fmt.Println("Server started on :8080")
		if err := server.ListenAndServe(); err != nil{
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	fmt.Println("\nShutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("failed to shutdown server: %v", err)
	}

	if err := s.Save(); err != nil {
		log.Printf("failed to save storage: %v", err)
	}

	fmt.Println("Server stopped.")
}
