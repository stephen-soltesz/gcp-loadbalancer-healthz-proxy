// gcp-loadbalancer-proxy provides an HTTP-to-insecure-HTTPS proxy.
//
// For example, this can be helpful when deploying GCP LoadBalancer health
// checks, which currently only support HTTP targets.
package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/m-lab/go/rtx"
)

var (
	rawURL  string
	rawAddr string

	// Create a ctx & cancel for
	ctx, cancelCtx = context.WithCancel(context.Background())
)

func init() {
	flag.StringVar(&rawURL, "url", "https://localhost:6443",
		"The base URL to query (insecurely).")
	flag.StringVar(&rawAddr, "port", ":8080",
		"Listen on the given address.")
}

func newReverseProxy(target *url.URL) (*httputil.ReverseProxy, error) {
	if target.RawQuery != "" || target.Path != "" {
		return nil, fmt.Errorf("Do not provide target path or queries. " +
			"All paths and queries are copied from incoming requests.")
	}
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
	}, nil
}

func main() {
	flag.Parse()

	target, err := url.Parse(rawURL)
	rtx.Must(err, "Failed to parse given url: %s", rawURL)

	reverseProxy, err := newReverseProxy(target)
	rtx.Must(err, "Failed to create reverse proxy: %s", target)

	http.Handle("/", reverseProxy)
	go func() {
		log.Fatal(http.ListenAndServe(rawAddr, nil))
	}()

	<-ctx.Done()
}
