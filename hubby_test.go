package hubby

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ClientMock struct{}

func TestGetDomains(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Basic c3RpYW46c3RpYW4=", r.Header.Get("Authorization"))
		w.Write([]byte(okResponse))
	})

	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	api := New("stian", "stian", httpClient)

	domains, err := api.GetDomains()
	if err != nil {
		fmt.Println(err)
	}

	assert.Nil(t, err)
	assert.Equal(t, 1, len(domains))

}

func TestDeleteDNSRecord(t *testing.T) {

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Basic c3RpYW46c3RpYW4=", r.Header.Get("Authorization"))
		w.WriteHeader(204)
		w.Write([]byte(""))
	})

	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	api := New("stian", "stian", httpClient)
	err := api.DeleteDNSRecord(99, 99)
	assert.Nil(t, err)

}

func TestUpdateDNSRecord(t *testing.T) {

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Basic c3RpYW46c3RpYW4=", r.Header.Get("Authorization"))
		w.WriteHeader(204)
		w.Write([]byte(""))
	})

	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	payload := DNSRecord{
		Host: "",
		TTL:  3600,
		Type: "A",
		Data: "1.1.1.1",
	}
	api := New("stian", "stian", httpClient)
	err := api.UpdateDNSRecord(99, 99, payload)
	assert.Nil(t, err)

}

//need some love..
const (
	okResponse = `[
		{"id": 1, "domain": "lus.re", "status":"active"}
		]`
)

func testingHTTPClient(handler http.Handler) (*http.Client, func()) {
	s := httptest.NewTLSServer(handler)

	cli := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, s.Listener.Addr().String())
			},
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	return cli, s.Close
}
