package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-microservices/handlers"
)

func main() {
	logger := log.New(os.Stdout, "[product-api]", log.LstdFlags)
	helloHandler := handlers.NewHello(logger)
	goodbyeHandler := handlers.NewGoodbye(logger)

	serveMux := http.NewServeMux()
	serveMux.Handle("/", helloHandler)
	serveMux.Handle("/bye", goodbyeHandler)

	http.ListenAndServe("localhost:9090", serveMux)
}
