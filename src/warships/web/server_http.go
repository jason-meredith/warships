package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"warships/net"
)

/*********************************************************
 *														 *
 *                   	  Warships						 *
 *					   Jason Meredith					 *
 *														 *
 *	DATE:		October 27, 2018						 *
 *	FILE: 		server_http.go							 *
 *	PURPOSE:	Runs a simple HTTP server to demo		 *
 *				the multithreading capabilities of Go	 *
 *				as well as a simple implementation of	 *
 *				the channels mechanism					 *
 *				 										 *
 *														 *
 *********************************************************/

// readFile reads a file into a string, in this case it reads html files
func readFile(filename string) string {
	data, err := ioutil.ReadFile(filename)

	if(err != nil) {
		panic(err)
	}

	return string(data)
}

// main Sets up and runs the web server. After integrated
func main() {

	// In the main program, the Server and message chan would be passed to the
	// http_server, not instantiated here
	server := new(net.Server)
	messages := make(chan string)

	// NOTE: This would be called by the function that creates and sets up the Game Server
	// This creates a new Goroutine, starting the runHttpServer() function in a new thread
	// I pass it the messages channel to pass messages back to this thread
	go runHttpServer(server, messages)

	for {
		// NOTE: Messages would be processed here, not necessarily simple printed to screen
		// Because I am outputting a channel here, it will block this thread until it receives
		// an input
		fmt.Println(<- messages)
	}


}

// runHttpServer runs a simple HTTP server, passing the Server to perform administrative tasks on
// as well as the message channel to send messages.
//
// Parameter message is a string channel, meaning that in this Goroutine, anytime I input (<-) a value into
// message, the thread blocks until the thread that created it reaches a point where it outputs the
// channel (line 54) where whatever was passed into the message channel here is outputted there
func runHttpServer(server *net.Server, message chan string ) {

	// When the user requests absolute path /style.css (really only requested by pages after loading)
	// return the stylesheet
	http.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		w.Write([]byte(readFile("html/css/style.css")))
		message <- "Incoming request: /style.css"
	})

	// Route /data will return JSON data about the server
	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{ json }"))
		message <- "Incoming request: /data"
	})

	// Route /about returns the about page
	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(readFile("html/about.html")))
		message <- "Incoming request: /about"
	})

	// Route / returns the root page showing server status info
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(readFile("html/index.html")))
		message <- "Incoming request: /"
	})

	// Show that the server is listening on port 8080 and start the server on that port
	fmt.Println("Exercise 4 by Jason Meredith")
	fmt.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
