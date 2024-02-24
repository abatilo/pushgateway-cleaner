package main

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func Test_FetchPushTimeMetric(t *testing.T) {
	mockHTTPClient := &MockHTTPClient{
		GetFunc: func(url string) (*http.Response, error) {
			mockResponse := `
# HELP push_time_seconds Last Unix time when changing this group in the Pushgateway succeeded.
# TYPE push_time_seconds gauge
push_time_seconds{instance="some_instance",job="some_job"} 1.7088061457877028e+09
`
			return &http.Response{
				Body: io.NopCloser(strings.NewReader(mockResponse)),
			}, nil
		},
	}

	_, err := fetchPushTimeMetric(mockHTTPClient, "http://localhost:9091")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
