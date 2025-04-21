package main

import (
	"log"
	"net/http"
)

func main() {
	//making a newserveMux
	const port = ":8080"
	serveMux := http.NewServeMux()

	//Handle(function): defines a specific url path on the server, associates path with handler
	//handler: responsible for processing incoming http requests and returning appropriate responses
	//fileserver: built in handler using http library, serves static files from specific directories on computer,
	//needs to know what directory to serve files from(used "." since that is current directory)
	serveMux.Handle("/", http.FileServer(http.Dir(".")))

	//making the server struct
	myServer := &http.Server{
		Addr:    port,
		Handler: serveMux,
	}

	//start an http server with the port and handler we created above/ handles any errors
	err := myServer.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Printf("Error: %v\n", err)
		log.Fatal("server error: ", err)
	}

}
