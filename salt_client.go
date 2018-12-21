package main

import (
	"fmt"
	"net/http"
)

type saltClient struct {
	Connector *saltConnector
}

func newSaltClient(config botConfig) (*saltClient, error) {
	c := saltClient{}
	c.Connector = newSaltConnector(config)
	err := c.Connector.auth()

	return &c, err
}

// executeJob method for salt client.
// this method is synchronous and will wait for job to finish and return result from execution
func (c *saltClient) executeJob(job saltJobPostData) (*http.Response, error) {
	var jobData string

	// If no arguments are specified, client will not sent arg parameter.
	// If arg is sent emptu or with whitespaces, will cause error on salt-api
	if job.Arguments == "" {
		jobData = fmt.Sprintf(`{"client": "local", "fun": "%s", "tgt": "%s", "kwarg": "%s"}`, job.Function, job.Target, job.Kwargs)
	} else {
		jobData = fmt.Sprintf(`{"client": "local", "fun": "%s", "arg": "%s", "tgt": "%s", "kwarg": "%s"}`, job.Function, job.Arguments, job.Target, job.Kwargs)
	}

	resp, err := c.Connector.post("", []byte(jobData))
	if err != nil {
		fmt.Println("Error while sending job execution:", err)
	}
	return resp, nil

}

// sentJob method for salt client.
// this method is asynchronous and will return only job id not result
func (c *saltClient) sentJob(job saltJobPostData) (*http.Response, error) {
	var jobData string

	// If no arguments are specified, client will not sent arg parameter.
	// If arg is sent emptu or with whitespaces, will cause error on salt-api
	if job.Arguments == "" {
		jobData = fmt.Sprintf(`{"client": "local", "fun": "%s", "tgt": "%s", "kwarg": ""}`, job.Function, job.Target)
	} else {
		jobData = fmt.Sprintf(`{"client": "local", "fun": "%s", "arg": "%s", "tgt": "%s", "kwarg": ""}`, job.Function, job.Arguments, job.Target)
	}

	resp, err := c.Connector.post("/minions", []byte(jobData))
	if err != nil {
		fmt.Println("Error while sending job to minions:", err)
	}
	return resp, nil

}

//get minions info
func (c *saltClient) getMinionsInfo() (*http.Response, error) {
	resp, err := c.Connector.get("/minions")
	if err != nil {
		fmt.Println("Error while getting minions info:", err)
		return nil, err
	}
	return resp, nil

}
