[![Go Report Card](https://goreportcard.com/badge/github.com/stephen-soltesz/gcp-loadbalancer-healthz-proxy)](https://goreportcard.com/report/github.com/stephen-soltesz/gcp-loadbalancer-healthz-proxy) [![Build Status](https://travis-ci.org/stephen-soltesz/gcp-loadbalancer-healthz-proxy.svg?branch=master)](https://travis-ci.org/stephen-soltesz/gcp-loadbalancer-healthz-proxy) [![Coverage Status](https://coveralls.io/repos/github/stephen-soltesz/gcp-loadbalancer-healthz-proxy/badge.svg?branch=master)](https://coveralls.io/github/stephen-soltesz/gcp-loadbalancer-healthz-proxy?branch=master)

# gcp-loadbalancer-proxy

gcp-loadbalancer-proxy provides an HTTP-to-insecure-HTTPS proxy.

For example, this can be helpful when deploying GCP LoadBalancer health
checks, which currently only support HTTP targets.

Default options support use as a proxy for HA Kubernetes master health
checks.