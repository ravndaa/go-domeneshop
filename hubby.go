package hubby

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

//Hubby ...
type Hubby struct {
	client *myhttp
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

//New domeneshop client.
func New(clientid string, clientsecret string, client *http.Client) *Hubby {
	if client == nil {
		client = &http.Client{}
	}

	apiclient := &myhttp{
		client: client,
		auth:   basicAuth(clientid, clientsecret),
	}

	api := Hubby{
		client: apiclient,
	}
	return &api
}

//GetDomains ...
func (a *Hubby) GetDomains() ([]Domain, error) {

	resp, err := a.client.GET("/domains")
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
	// make it cleaner net/url pacjage ?
	path := fmt.Sprintf("/domains/%v/dns", domainid)
	if dnstype != "" && host != "" {
		path = fmt.Sprintf("%v?host=%v&type=%v", path, host, dnstype)
	} else if host != "" {
		path = fmt.Sprintf("%v?host=%v", path, host)
	} else if dnstype != "" {
		path = fmt.Sprintf("%v?type=%v", path, dnstype)
	}
	resp, err := a.client.GET(path)
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

// Validate DNSRecord before sending it.
func validateDNSRecord(record DNSRecord) bool {
	if record.Host == "" {
		return false
	}
	if record.Type == "" {
		return false
	}
	if record.Data == "" {
		return false
	}

	switch DNSType := record.Type; DNSType {
	case "MX":
		if record.Priority == "" {
			return false
		}
	case "SRV":
		if record.Priority == "" {
			return false
		}
		if record.Weight == "" {
			return false
		}
		if record.Port == "" {
			return false
		}
	}
	return true
}

//AddDNSRecord ...
func (a *Hubby) AddDNSRecord(domainid int, value DNSRecord) error {

	if validateDNSRecord(value) == false {
		return errors.New(ErrMissingRequiredField)
	}

	path := fmt.Sprintf("/domains/%v/dns", domainid)
	resp, err := a.client.POST(path, value)
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

	url := fmt.Sprintf("/domains/%v/dns/%v", domainid, dnsrecordid)
	resp, err := a.client.PUT(url, value)
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
	url := fmt.Sprintf("/domains/%v/dns/%v", domainid, dnsrecordid)

	resp, err := a.client.DELETE(url)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 300 {
		return errors.New(ErrNotSureYet)
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
	DNS       bool   `json:"dns"`
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

const (
	//ErrMissingRequiredField used for checking fields required.
	ErrMissingRequiredField = "missing required field"
	//ErrNotSureYet ...
	ErrNotSureYet = "not sure what happend."
)
