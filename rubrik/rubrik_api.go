//
// rubrik-exporter
//
// Exports metrics from rubrik backup for prometheus
//
// License: Apache License Version 2.0,
// Organization: Claranet GmbH
// Author: Martin Weber <martin.weber@de.clara.net>
//

package rubrik

import (
	//	"os"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type RequestParams struct {
	body, header string
	params       url.Values
}

type Rubrik struct {
	url      string
	username string
	password string

	sessionToken string
	isLoggedIn   bool
}

func (r *Rubrik) makeRequest(reqType string, action string, p RequestParams) (*http.Response, error) {
	log.Printf("Is logged in: %t", r.isLoggedIn)

	_url := r.url + action

	log.Printf("Requested action: %s", action)

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	var netClient = http.Client{Transport: tr}

	body := p.body

	_url += "?" + p.params.Encode()
	log.Printf("Request full URL: %s", _url)

	req, err := http.NewRequest(reqType, _url, strings.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "text/JSON")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", r.sessionToken))

	resp, err := netClient.Do(req)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// Check HTTP status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("API Error: HTTP %d from %s", resp.StatusCode, action)
		resp.Body.Close()
		// Return empty response with error to prevent JSON parsing of error pages
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, action)
	}

	return resp, nil
}

// NewRubrik - Creates a new Rubrik API instance and login to it
func NewRubrik(url string, username string, password string) *Rubrik {

	log.Print("Create new API Instance")
	session := &Rubrik{
		url:          url,
		username:     username,
		password:     password,
		sessionToken: "",
		isLoggedIn:   false,
	}
	session.Login()
	log.Printf("Session-Token: %s", session.sessionToken)

	return session
}
