package net

import (
	"net/http"
	"time"
)

//NewHttpClient is constructor for http.Client 
func NewHttpClient(requestTimeout time.Duration, checkRedirect func(req *http.Request, via []*http.Request) error) *http.Client {
	return &http.Client{
		Timeout:       requestTimeout,
		CheckRedirect: checkRedirect,
	}
}
