package main

import "net/http"

type HTTPClient interface {
	Get(url string) (*http.Response, error)
	Do(req *http.Request) (*http.Response, error)
}

type MockHTTPClient struct {
	GetFunc func(url string) (*http.Response, error)
	DoFunc  func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Get(url string) (*http.Response, error) {
	return m.GetFunc(url)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}
