package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-microservices/handlers"
)

func main() {
	logger := log.New(os.Stdout, "[product-api]", log.LstdFlags)
	helloHandler := handlers.NewHello(logger)
	goodbyeHandler := handlers.NewGoodbye(logger)

	serveMux := http.NewServeMux()
	serveMux.Handle("/", helloHandler)
	serveMux.Handle("/bye", goodbyeHandler)

	server := &http.Server{
		Addr: ":9090",
		Handler: serveMux,
		IdleTimeout: 120 * time.Second,
		ReadTimeout: 1 * time.Second,
		WriteTimeout: 1 * time.Second,
	} 
	
	go func() {
		err := server.ListenAndServe()

		if err != nil {
			logger.Fatal(err)
		}
	}()
	
	sigchan := make(chan os.Signal)
	signal.Notify(sigchan, os.Interrupt)
	signal.Notify(sigchan, os.Kill)

	logger.Println("Received terminate, graceful shutdown", <- sigchan)

	timeoutContext, _ := context.WithTimeout(context.Background(), 30 * time.Second)
	server.Shutdown(timeoutContext)
}
