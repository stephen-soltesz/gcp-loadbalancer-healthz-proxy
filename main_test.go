// gcp-service-discovery provides an HTTP-to-insecure-HTTPS proxy.
//
// For example, this can be helpful for Google LoadBalancer health
// checkers which only support HTTP health-check targets.
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"testing"

	"github.com/m-lab/go/rtx"
)

func Test_newReverseProxy(t *testing.T) {
	tests := []struct {
		name        string
		target      *url.URL
		statusCode  int
		fileContent string
		path        string
		want        *httputil.ReverseProxy
		wantErr     bool
		targetPath  string
	}{
		{
			name:        "success",
			statusCode:  http.StatusOK,
			path:        "/",
			fileContent: "ok",
		},
		{
			name:        "success-status-code-propagates",
			statusCode:  http.StatusNotFound,
			path:        "/",
			fileContent: "ok",
		},
		{
			name:        "success-path-propagates",
			statusCode:  http.StatusOK,
			path:        "/healthz",
			fileContent: "ok",
		},
		{
			name:        "success-path-with-query-parameters",
			statusCode:  http.StatusOK,
			path:        "/healthz?foo=bar&thing=that",
			fileContent: "ok",
		},
		{
			name:        "failure-reverse-proxy-returns-error",
			statusCode:  http.StatusOK,
			path:        "/",
			fileContent: "ok",
			targetPath:  "/bar",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a tls target server that the reverse proxy will query.
			tls := httptest.NewTLSServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tt.statusCode)
					fmt.Fprint(w, tt.fileContent+r.URL.String())
				}),
			)
			defer tls.Close()

			// Setup the reverse proxy & test server to target the tls test server.
			u, err := url.Parse(tls.URL + tt.targetPath)
			rtx.Must(err, "httptest URL could not be parsed: %s", tls.URL)
			rp, err := newReverseProxy(u)
			if err != nil {
				if tt.targetPath != "" {
					// If the target path is non zero, then we want this to be an error.
					return
				} else {
					t.Errorf("newReverseProxy error = %v, targetPath %v", err, tt.targetPath)
				}
			}
			ts := httptest.NewServer(rp)
			defer ts.Close()

			// Issue a request to the second test server and verify return values.
			resp, err := http.Get(ts.URL + tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("newReverseProxy request error = %v, wantErr %v", err, tt.wantErr)
			}
			if resp.StatusCode != tt.statusCode {
				t.Errorf("newReverseProxy request StatusCode; got = %v, want %v",
					resp.StatusCode, tt.statusCode)
			}
			b, err := ioutil.ReadAll(resp.Body)
			rtx.Must(err, "Failed to read response")
			if string(b) != tt.fileContent+tt.path {
				t.Errorf("newReverseProxy fileContent; got = %v, want %v",
					string(b), tt.fileContent+tt.path)
			}

		})
	}
}

func Test_main(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "ok")
		}),
	)
	defer ts.Close()
	rawURL = ts.URL
	go func() {
		http.Get("http://localhost" + rawAddr)
		cancelCtx()
	}()
	main()
}
