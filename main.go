// gcp-service-discovery provides an HTTP-to-insecure-HTTPS proxy.
//
// For example, this can be helpful for Google LoadBalancer health
// checkers which only support HTTP health-check targets.
package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/m-lab/go/rtx"
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

func newReverseProxy(target *url.URL) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		// Override the http "Host:" header to specify target.Host.
		req.Host = target.Host
		// Override the request scheme and host to connect via TLS.
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
	}

	// Override the default transport settings to support insecure TLS connections.
	transport := http.DefaultTransport.(*http.Transport)
	transport.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	return &httputil.ReverseProxy{
		Director:  director,
		Transport: transport,
	}
}

func main() {
	flag.Parse()

	target, err := url.Parse(rawURL)
	rtx.Must(err, "Failed to parse given url: %s", rawURL)
	if target.RawQuery != "" || target.RawPath != "" {
		log.Fatal("Do not provide target path or queries. " +
			"All paths and queries are copied from incoming requests.")
	}

	http.Handle("/", newReverseProxy(target))
	log.Fatal(http.ListenAndServe(rawAddr, nil))
}
