package hubby

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	baseURL = "https://api.domeneshop.no/v0"
)

//Hubby ...
type Hubby struct {
	client *http.Client
	auth   string
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

//New domeneshop client.
func New(clientid string, clientsecret string, client *http.Client) *Hubby {

	api := Hubby{
		auth:   basicAuth(clientid, clientsecret),
		client: client,
	}
	return &api
}

//GetDomains ...
func (a *Hubby) GetDomains() ([]Domain, error) {
	url := fmt.Sprintf("%v/domains", baseURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "ravndaa/hubby")
	req.Header.Add("Authorization", "Basic "+a.auth)
	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, err
	}

	domains := []Domain{}
	err = json.Unmarshal(body, &domains)
	if err != nil {
		return nil, err
	}

	return domains, nil

}

//GetDomainDNSRecords ...
func (a *Hubby) GetDomainDNSRecords(domainid int) ([]DNSRecord, error) {
	url := fmt.Sprintf("%v/domains/%v/dns", baseURL, domainid)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "ravndaa/hubby")
	req.Header.Add("Authorization", "Basic "+a.auth)
	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, err
	}
	records := []DNSRecord{}
	err = json.Unmarshal(body, &records)
	if err != nil {
		return nil, err
	}
	return records, nil
}

//AddDNSRecord ...
func (a *Hubby) AddDNSRecord(domainid int, value DNSRecord) error {
	url := fmt.Sprintf("%v/domains/%v/dns", baseURL, domainid)

	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(value)

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "ravndaa/hubby")
	req.Header.Add("Authorization", "Basic "+a.auth)
	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 201 {
		fmt.Println(resp)

	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	return nil

}

//UpdateDNSRecord ...
func (a *Hubby) UpdateDNSRecord(domainid int, dnsrecordid int, value DNSRecord) error {
	url := fmt.Sprintf("%v/domains/%v/dns/%v", baseURL, domainid, dnsrecordid)

	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(value)

	req, err := http.NewRequest("PUT", url, payload)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "ravndaa/hubby")
	req.Header.Add("Authorization", "Basic "+a.auth)
	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return err
	}
	if resp.StatusCode != 204 {

		return fmt.Errorf("%s", body)
	}
	return nil
}

//DeleteDNSRecord ...
func (a *Hubby) DeleteDNSRecord(domainid int, dnsrecordid int) error {
	url := fmt.Sprintf("%v/domains/%v/dns/%v", baseURL, domainid, dnsrecordid)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "ravndaa/hubby")
	req.Header.Add("Authorization", "Basic "+a.auth)
	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 204 {
		return errors.New("noe gikk galt")
	}

	return nil
}

//Domain ..
type Domain struct {
	ID             int    `json:"id"`
	Domain         string `json:"domain"`
	Status         string `json:"status"`
	ExpiryDate     string `json:"expiry_date"`
	RegisteredDate string `json:"registered_date"`
	Renew          string `json:"renew"`
	Registrant     string `json:"registrant"`
}

//DNSRecord ...
type DNSRecord struct {
	ID   int    `json:"id,omitempty"`
	Host string `json:"host"`
	TTL  int    `json:"ttl"`
	Type string `json:"type"`
	Data string `json:"data"`
}
