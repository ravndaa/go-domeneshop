package domeneshop

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

/*

Basic testing seems to work, but need to add more tests and error handling.

*/

func TestGetDomains(t *testing.T) {

	okResponse := `[
		{"id": 1, "domain": "lus.re", "status":"active"},
		{"id": 2, "domain": "norge.no", "status":"active"}
		]`
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Basic c3RpYW46c3RpYW4=", r.Header.Get("Authorization"))
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/v0/domains", r.RequestURI)
		_, err := w.Write([]byte(okResponse))
		if err != nil {
			fmt.Println(err)
		}
	})

	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	api, _ := New("stian", "stian", httpClient)

	domains, err := api.GetDomains("")
	if err != nil {
		fmt.Println(err)
	}

	assert.Nil(t, err)
	assert.Equal(t, 2, len(domains))
	assert.Equal(t, "lus.re", domains[0].Domain)

}

func TestGetDomainsWithFilter(t *testing.T) {
	okResponse := `[
		{"id": 2, "domain": "norge.no", "status":"active"}
		]`
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Basic c3RpYW46c3RpYW4=", r.Header.Get("Authorization"))
		assert.Equal(t, "/v0/domains?domain=.no", r.RequestURI)
		assert.Equal(t, "GET", r.Method)
		_, err := w.Write([]byte(okResponse))
		if err != nil {
			fmt.Println(err)
		}
	})

	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	api, _ := New("stian", "stian", httpClient)

	domains, err := api.GetDomains(".no")
	if err != nil {
		fmt.Println(err)
	}

	assert.Nil(t, err)
	assert.Equal(t, "norge.no", domains[0].Domain)

}

func TestDeleteDNSRecord(t *testing.T) {

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Basic c3RpYW46c3RpYW4=", r.Header.Get("Authorization"))
		assert.Equal(t, "DELETE", r.Method)
		w.WriteHeader(204)
		_, err := w.Write([]byte(""))
		if err != nil {
			fmt.Println(err)
		}
	})

	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	api, _ := New("stian", "stian", httpClient)
	err := api.DeleteDNSRecord(99, 99)
	assert.Nil(t, err)

}

func TestUpdateDNSRecord(t *testing.T) {

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Basic c3RpYW46c3RpYW4=", r.Header.Get("Authorization"))
		assert.Equal(t, "PUT", r.Method)
		w.WriteHeader(204)
		_, err := w.Write([]byte(""))
		if err != nil {
			fmt.Println(err)
		}
	})

	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	payload := DNSRecord{
		Host: "",
		TTL:  3600,
		Type: "A",
		Data: "1.1.1.1",
	}
	api, _ := New("stian", "stian", httpClient)
	err := api.UpdateDNSRecord(99, 99, payload)
	assert.Nil(t, err)

}

//
func TestListDNSRecord(t *testing.T) {
	okResponse := `[
{
"id": 1,
"host": "@",
"ttl": 3600,
"type": "A",
"data": "192.168.0.1"
}
		]`

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Basic c3RpYW46c3RpYW4=", r.Header.Get("Authorization"))
		assert.Equal(t, "GET", r.Method)
		w.WriteHeader(200)
		_, err := w.Write([]byte(okResponse))
		if err != nil {
			fmt.Println(err)
		}
	})

	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	api, _ := New("stian", "stian", httpClient)
	records, err := api.ListDNSRecords(1, "", "")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(records))
}

//
func TestListDNSRecordWithFilters(t *testing.T) {
	okResponse := `[
{
"id": 1,
"host": "@",
"ttl": 3600,
"type": "A",
"data": "192.168.0.1"
}
		]`

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Basic c3RpYW46c3RpYW4=", r.Header.Get("Authorization"))
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/v0/domains/1/dns?host=norge.no&type=A", r.RequestURI)
		w.WriteHeader(200)
		_, err := w.Write([]byte(okResponse))
		if err != nil {
			fmt.Println(err)
		}
	})

	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	api, _ := New("stian", "stian", httpClient)
	records, err := api.ListDNSRecords(1, "norge.no", "A")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(records))
}

func TestAddDNSRecord(t *testing.T) {
	okResponse := `[
{
"id": 1,
"host": "@",
"ttl": 3600,
"type": "A",
"data": "192.168.0.1"
}
		]`

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Basic c3RpYW46c3RpYW4=", r.Header.Get("Authorization"))
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/v0/domains/1/dns", r.RequestURI)
		w.WriteHeader(201)
		_, err := w.Write([]byte(okResponse))
		if err != nil {
			fmt.Println(err)
		}
	})

	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	record := DNSRecord{Host: "test", Type: "A", Data: "127.0.0.1"}
	api, _ := New("stian", "stian", httpClient)
	err := api.AddDNSRecord(1, record)
	assert.Nil(t, err)
}

// Helpers

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
