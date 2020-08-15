package main

import (
	"log"
	"mpc_sample_project/server"
	"net/http"
)

func main() {
	srv := server.CreateServer()
	err := srv.ListenAndServe()
	if err == http.ErrServerClosed{
		log.Println("server shuting down gracefully")
	} else {
		log.Println("unexpected server shutdown...")
		log.Println("ERR: ", err)
	}
}
