package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type saltConnector struct {
	Config    botConfig
	Client    *http.Client
	AuthToken string
}

func newSaltConnector(config botConfig) *saltConnector {
	c := saltConnector{Config: config}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: config.Salt.SSLSkipVerify,
		},
	}
	c.Client = &http.Client{Transport: tr}

	return &c
}

// auth method for saltConnector to get auth token with limited time validity
func (c *saltConnector) auth() error {
	conf := c.Config.Salt
	authURL := fmt.Sprintf("https://%s:%d/login", conf.URL, conf.Port)
	authData := fmt.Sprintf(`{ "username":"%s", "password":"%s", "eauth": "%s" }`, conf.User, conf.Password, conf.Eauth)

	req, err := http.NewRequest("POST", authURL, bytes.NewBuffer([]byte(authData)))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := c.Client.Do(req)
	if err != nil {
		fmt.Printf("Failed to authentificate:%v", err)
		return err
	}

	if resp.StatusCode != 200 {
		return err
	}

	c.AuthToken = resp.Header.Get("X-Auth-Token")
	return nil
}

// post method for saltConnector
func (c *saltConnector) post(endpoint string, data []byte) (*http.Response, error) {
	conf := c.Config.Salt
	var url string

	// Possible to specify another than defualt end point for salt api
	// Official documentation: https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html
	// If no endpoint is specified than defualt / is used
	if len(strings.TrimSpace(endpoint)) == 0 {
		url = fmt.Sprintf("https://%s:%d/", conf.URL, conf.Port)
	} else {
		url = fmt.Sprintf("https://%s:%d%s", conf.URL, conf.Port, endpoint)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Token", c.AuthToken)

	resp, err := c.Client.Do(req)
	if err != nil {
		fmt.Printf("Client POST failed: %v", err)
		return nil, err
	}
	if resp == nil {
		fmt.Printf("Client POST failed with nil: %v", err)
		return nil, err
	}

	// Check for response code. all 2** are ok, others not
	if match, _ := regexp.MatchString("2[0-9]*", string(resp.StatusCode)); match {
		fmt.Printf("Failed to POST request!. Ended with code: %v. Error: %v", resp.StatusCode, err)
		return nil, err
	}

	return resp, nil
}

// get method for saltConnector
func (c *saltConnector) get(endpoint string) (*http.Response, error) {
	conf := c.Config.Salt
	var url string

	// Possible to specify another than defualt end point for salt api
	// Official documentation: https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html
	// If no endpoint is specified than defualt / is used
	if len(strings.TrimSpace(endpoint)) == 0 {
		url = fmt.Sprintf("https://%s:%d/", conf.URL, conf.Port)
	} else {
		url = fmt.Sprintf("https://%s:%d%s", conf.URL, conf.Port, endpoint)
	}

	// fmt.Println("URL formated:", url)
	fmt.Println("URL2 formated:", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Token", c.AuthToken)

	resp, err := c.Client.Do(req)
	if err != nil {
		fmt.Printf("Client GET failed: %v", err)
		return nil, err
	}
	if resp == nil {
		fmt.Printf("Client GET failed with nil: %v", err)
		return nil, err
	}

	// Check for response code. all 2** are ok, others not
	if match, _ := regexp.MatchString("2[0-9]*", string(resp.StatusCode)); match {
		fmt.Printf("Failed to POST request!. Ended with code: %v. Error: %v", resp.StatusCode, err)
		return nil, err
	}

	return resp, nil
}
