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

//FindDomains with filter
func (a *Hubby) FindDomains(filter string) ([]Domain, error) { return nil, nil }

//FindDomain using id
func (a *Hubby) FindDomain(domainid string) ([]Domain, error) { return nil, nil }

//ListDNSRecords ...
func (a *Hubby) ListDNSRecords(domainid int, host string, dnstype string) ([]DNSRecord, error) {
	url := fmt.Sprintf("%v/domains/%v/dns", baseURL, domainid)
	// add some queries if not empty
	if dnstype != "" && host != "" {
		url = fmt.Sprintf("%v?host=%v&type=%v", url, host, dnstype)
	} else if host != "" {
		url = fmt.Sprintf("%v?host=%v", url, host)
	} else if dnstype != "" {
		url = fmt.Sprintf("%v?typ=%v", url, dnstype)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	// Set some request header data, one required and one for the fun.
	req.Header.Set("User-Agent", "ravndaa/hubby")
	req.Header.Add("Authorization", "Basic "+a.auth)

	// ask the api nicely to get some dns records.
	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	// check if return body is nil, or move on.
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	// read the response body since it isnt nil.
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, err
	}
	// create an Dnsrecord array used by json unmarshal.
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

	if value.Host == "" {
		return errors.New("missing host")
	}
	if value.Type == "" {
		return errors.New("mssing type")
	}
	if value.Data == "" {
		return errors.New("missing data")
	}

	switch DNSType := value.Type; DNSType {
	case "MX":
		if value.Priority == "" {
			return errors.New("missing priority")
		}
	case "SRV":
		if value.Priority == "" {
			return errors.New("missing priority")
		}
		if value.Weight == "" {
			return errors.New("missing weight")
		}
		if value.Port == "" {
			return errors.New("missing port")
		}
	}

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
	// not sure why I have copied this in, should check it some time.
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return err
	}

	if resp.StatusCode >= 300 {
		return fmt.Errorf("%s", body)
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
	ID             int      `json:"id"`
	Domain         string   `json:"domain"`
	Status         string   `json:"status"`
	ExpiryDate     string   `json:"expiry_date"`
	RegisteredDate string   `json:"registered_date"`
	Renew          bool     `json:"renew"`
	Registrant     string   `json:"registrant"`
	Nameservers    []string `json:"nameservers"`
	Services       Service  `json:"services"`
}

//Service ...
type Service struct {
	Registrar bool   `json:"registrar"`
	Dns       bool   `json:"dns"`
	Email     bool   `json:"email"`
	Webhotel  string `json:"webhotel"`
}

//DNSRecord ...
type DNSRecord struct {
	ID       int    `json:"id,omitempty"`
	Host     string `json:"host,omitempty"`
	TTL      int    `json:"ttl,omitempty"`
	Type     string `json:"type,omitempty"`
	Data     string `json:"data,omitempty"`
	Priority string `json:"priority,omitempty"`
	Weight   string `json:"weight,omitempty"`
	Port     string `json:"port,omitempty"`
}
