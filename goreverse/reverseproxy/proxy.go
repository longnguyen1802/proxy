package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

func AddCookie(rw *http.ResponseWriter, cookie *http.Cookie) {
	http.SetCookie(*rw, cookie)
}

func reverseString(input string) string {
	runes := []rune(input)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func SetHeader(destination *http.Header, source *http.Header) {
	for key, values := range *source {
		for _, value := range values {
			destination.Add(key, value)
		}
	}
}
func EncryptCookie(cookie *http.Cookie) {
	cookie.Value = reverseString(cookie.Value)
}

func DecryptCookie(cookie *http.Cookie) {
	cookie.Value = reverseString(cookie.Value)
}

func EncryptResponseCookie(rw *http.ResponseWriter, cookies []*http.Cookie) {
	for _, cookie := range cookies {
		EncryptCookie(cookie)
		AddCookie(rw, cookie)
	}
}
func DecryptRequestCookie(req *http.Request) {
	originalCookies := req.Cookies()
	req.Header.Del("Cookie")
	for _, cookie := range originalCookies {
		DecryptCookie(cookie)
		req.AddCookie(cookie)
	}
}

func SetCookie(rw *http.ResponseWriter, cookie *http.Cookie) {
	http.SetCookie(*rw, cookie)
}
func main() {
	// define origin server URL
	originServerURL, err := url.Parse("https://github.com")
	if err != nil {
		log.Fatal("invalid origin server URL")
	}

	reverseProxy := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		fmt.Printf("[reverse proxy server] received request at: %s\n", time.Now())

		// set req Host, URL and Request URI to forward a request to the origin server

		// Use a Goroutine to process the request and response

		req.Host = originServerURL.Host
		req.URL.Host = originServerURL.Host
		req.URL.Scheme = originServerURL.Scheme
		req.RequestURI = ""
		//DecryptRequestCookie(req)
		originServerResponse, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("Error while processing request: %s", err)
			return
		}
		for key, values := range originServerResponse.Header {
			// if key != "Set-Cookie" {
			// 	for _, value := range values {
			// 		rw.Header().Add(key, value)
			// 	}
			// }
			for _, value := range values {
				rw.Header().Add(key, value)
			}
		}
		//EncryptResponseCookie(&rw, originServerResponse.Cookies())
		// return response to the client
		rw.WriteHeader(originServerResponse.StatusCode)
		io.Copy(rw, originServerResponse.Body)
	})

	// Run the HTTP server in a Goroutine

	go func() {
		log.Fatal(http.ListenAndServe(":8080", reverseProxy))
	}()

	done := make(chan struct{})
	<-done

}
