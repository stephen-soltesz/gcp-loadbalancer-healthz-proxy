package main

import (
	"crypto/tls"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

var (
	rawURL  string
	rawAddr string
)

func init() {
	flag.StringVar(&rawURL, "url", "https://localhost:6443",
		"The URL to query (insecurely).")
	flag.StringVar(&rawAddr, "port", ":8080",
		"Listen on the given address.")
}

func main() {
	flag.Parse()

	// A client for making all backend requests.
	c := http.Client{
		Transport: &http.Transport{
			// Setup support for insecure TLS connections.
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			// Duplicate settings from http.DefaultTransport.
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Issue request, and read full response.
		resp, err := c.Get(rawURL + r.URL.Path)
		if err != nil {
			log.Printf("Error getting (%s%s): %v", rawURL, r.URL.Path, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// Report the same status and content.
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading body of (%s%s): %v", rawURL, r.URL.Path, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// Report the same status and content.
		// Note: we ignore headers.
		w.WriteHeader(resp.StatusCode)
		w.Write(body)
	})

	log.Fatal(http.ListenAndServe(rawAddr, nil))
}
