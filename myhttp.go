package hubby

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	baseURL = "https://api.domeneshop.no/v0"
)

type myhttp struct {
	client *http.Client
	auth   string
}

// needs /domains or simular
func (m myhttp) GET(path string) (*http.Response, error) {
	url := fmt.Sprintf("%v%v", baseURL, path)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "ravndaa/hubby")
	req.Header.Add("Authorization", "Basic "+m.auth)
	resp, err := m.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m myhttp) POST(path string, data interface{}) (*http.Response, error) {
	url := fmt.Sprintf("%v%v", baseURL, path)
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(data)

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "ravndaa/hubby")
	req.Header.Add("Authorization", "Basic "+m.auth)
	resp, err := m.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m myhttp) PUT(path string, data interface{}) (*http.Response, error) {
	url := fmt.Sprintf("%v%v", baseURL, path)

	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(data)

	req, err := http.NewRequest("PUT", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "ravndaa/hubby")
	req.Header.Add("Authorization", "Basic "+m.auth)
	resp, err := m.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m myhttp) DELETE(path string) (*http.Response, error) {
	url := fmt.Sprintf("%v%v", baseURL, path)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "ravndaa/hubby")
	req.Header.Add("Authorization", "Basic "+m.auth)
	resp, err := m.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
