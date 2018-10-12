package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func init() {

}

func main() {
	target, err := url.Parse("https://www.google.com/")
	if err != nil {
		log.Fatal(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(target)

	http.Handle("/", proxy)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
