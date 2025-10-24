package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type APIRequest struct {
	Body    string
	Headers map[string]string
	Method  string
	Query   map[string]string
	URL     url.URL
}

type APIResponse struct {
	Success    bool
	Error      error
	StatusCode int
	Status     string
	Body       []byte
	Headers    http.Header
}

func (ac *APIClient) Request(request APIRequest) (response APIResponse) {
	var (
		requestBody io.Reader = nil
	)

	response.Success = false

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSec*time.Second)
	defer cancel()

	if request.Body != "" {
		requestBody = strings.NewReader(request.Body)
	}

	req, err := http.NewRequestWithContext(ctx, request.Method, request.URL.String(), requestBody)
	if err != nil {
		response.Error = err
		return response
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Accept", "application/json")

	resp, err := ac.Client.Do(req)
	if err != nil {
		response.Error = err
		return response
	}
	defer resp.Body.Close()

	response.Status = resp.Status
	response.StatusCode = resp.StatusCode

	if resp.StatusCode != http.StatusOK {
		response.Error = fmt.Errorf("non-200 code: %d: %s", resp.StatusCode, resp.Status)
		return response
	}

	response.Body, err = io.ReadAll(resp.Body)
	if err != nil {
		response.Error = err
		return response
	}

	response.Success = true

	return response
}
