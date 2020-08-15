package server

import (
	"log"
	"net/http"
	"os"
	"os/signal"
)

func makeServer(addr string, r http.Handler) *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}

func handleGracefulShutdown(server *http.Server) {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	<-quit
	log.Println("receive interrupt signal")
	if err := server.Close(); err != nil {
		log.Fatal("Server Close:", err)
	}
}
