package net

import (
	"net/http"
	"time"
)

//NewHttpClient is constructor for http.Client 
func NewHttpClient(transport *http.Transport, requestTimeout time.Duration, checkRedirect func(req *http.Request, via []*http.Request) error) *http.Client {
	return &http.Client{
		Transport: transport,
		Timeout:       requestTimeout,
		CheckRedirect: checkRedirect,
	}
}
