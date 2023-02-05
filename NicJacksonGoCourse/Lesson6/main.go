package main

import (
	"Lesson3/handlers"
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	l := log.New(os.Stdout, "coffee-shop-api", log.LstdFlags)

	ph := handlers.NewProducts(l)

	sm := mux.NewRouter()

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", ph.GetProducts)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", ph.AddProduct)
	postRouter.Use(ph.MiddlewareProductValidation)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/id:[0-9]+", ph.UpdateProduct)
	putRouter.Use(ph.MiddlewareProductValidation)

	s := http.Server{
		Addr:         ":8080",
		Handler:      sm,
		IdleTimeout:  30 * time.Second,
		WriteTimeout: 1 * time.Second,
		ReadTimeout:  1 * time.Second,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			log.Fatalln("error starting server!")
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Received terminate, graceful shutdown....", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)

}
