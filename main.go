package main

import (
	"Practice/handlers"
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	l := log.New(os.Stdout, "Nurture", log.LstdFlags) // logging to os.stdout file, can also prefix the log, flag std
	ph := handlers.NewProducts(l)                     // creating reference to the handler and injecting l in that
	// register handler with server
	// http.HandleFunc(), we are converting the function into a handler type and then registering that into a default serve MUX (DSM)
	// sm := http.NewServeMux() // discard this, use gorilla service instead
	sm := mux.NewRouter() // create a new serve MUX
	// specified at this level that we only want to work with products
	// instead use a sub-router. Set the routing path
	getRouter := sm.Methods("GET").Subrouter() // call the get router (check docs)
	//above is best way to add method, not below one
	// sm.Handle("/products", ph)                 //products will act as the main get handler itself
	getRouter.HandleFunc("/", ph.GetProducts) // handles a func, we want it to directly handle the "get" function

	// create a sub-router for put requests
	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProducts) // define a varibale for gorilla mux and use regexp. regexp is within {}
	// regexp will be automatically extracted out of the URL by the router
	putRouter.Use(ph.MiddlewareProductValidation) // passing the middleware via the USE function

	//middleware will validate the objects passed in the method

	// create a sub-router for post requests
	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", ph.AddProduct)
	postRouter.Use(ph.MiddlewareProductValidation)
	s := &http.Server{ // set some tunable parameters for optimization
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Received terminate and graceful shutdown", sig)
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}
