package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	//making a newserveMux
	const port = ":8080"

	//keeps count of how many requests are being made
	var counter apiConfig

	//mux or multiplexer
	//it is a request router
	// it gets incoming http requests and decides which handler function should process the request
	//maps url patterns to handler functions
	serveMux := http.NewServeMux()

	//handlefunc register handlers with serveMux
	//takes in the "/healthz" endpoint
	//takes in a function with the signature "func(http.ResponseWriter, *http.Request)"
	//It automatically converts your function to a handler interface
	serveMux.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	//registers handler with serveMux
	//takes url path with an object with the method "ServeHTTP(http.ResponseWriter, *http.Request)"
	//used for pre-built handlers or custom handler type
	//want to use this over a handle in more complex situations such as with the fileserver handler or using miiddleware like stripPrefix
	//Strip prefix takes away the prefix "/app" from the handler
	//FileServer is a built in handler, automatically handles file serving, content types, and directory listings
	//FileServer serves static content
	appHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	serveMux.Handle("/app/", counter.MiddlewareMetricsInc(appHandler))

	//register the metrics handler
	serveMux.HandleFunc("/metrics", counter.RequestNum)

	//register the reset handler
	serveMux.HandleFunc("/reset", counter.resetNum)

	//making the server struct
	myServer := &http.Server{
		Addr:    port,
		Handler: serveMux,
	}

	//start an http server with the port and handler we created above/ handles any errors
	log.Println("Starting server on port #", port)
	err := myServer.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatal("server error: ", err)
	}

}
