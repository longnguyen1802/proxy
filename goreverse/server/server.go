package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func AddCookie(rw *http.ResponseWriter, cookie *http.Cookie) {
	http.SetCookie(*rw, cookie)
}

func main() {
	originServerHandler := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		fmt.Println("====================================================================")
		fmt.Printf("Request from client %+v", req)
		cookie := http.Cookie{
			Name:    "mycookie",
			Value:   "Hello, this is my cookie!",
			Expires: time.Now().Add(24 * time.Hour), // Cookie will expire in 24 hours
			Path:    "/",
		}

		// Set the cookie in the response
		AddCookie(&rw, &cookie)
		fmt.Printf("[origin server] received request at: %s\n", time.Now())
		_, _ = fmt.Fprint(rw, "origin server response\n")
	})

	log.Fatal(http.ListenAndServe(":8081", originServerHandler))
}
