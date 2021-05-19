package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"github.com/spf13/viper"
	"github.com/op/go-logging"
)

type Client struct {
	user string
	password string
	piot_host string
	log *logging.Logger
	token string
}

func NewClient(logger *logging.Logger) *Client {
	client := &Client{}
	client.log = logger
	client.user = viper.GetString("user")
	client.password = viper.GetString("password")
	client.piot_host = viper.GetString("piot.host")
	client.token = ""

	client.log.Debug("New instance of api client created:")
	client.log.Debugf("  user: %s", client.user)
	client.log.Debugf("  piot host: %s", client.piot_host)

	return client
}

func (c *Client) execute(method string, path string, body *[]byte) (*http.Response, error) {
	var bodyIoReader io.Reader

	var url string
	url = fmt.Sprintf("https://%s/%s", c.piot_host, path)

	c.log.Debugf("%s Request to: %s", method, url)
	if body != nil {
		c.log.Debugf("Request body: %s", string(*body))
		bodyIoReader = bytes.NewBuffer(*body)
	}

	req, err := http.NewRequest(method, url, bodyIoReader)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")

	// Add this header only for requests with body
	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	if (c.token == "") {
		c.log.Debug("Setting basic authorization (header)")
		req.SetBasicAuth(c.user, c.password)
	} else {
		c.log.Debugf("Setting bearer authorization (reusing token %s)", c.token)
		req.Header.Add("Authorization", "Bearer " + c.token)
	}
	httpClient := &http.Client{}

	resp, err := httpClient.Do(req)
	if err != nil {
		return resp, err
	}

	c.log.Debugf("Response Status: %s", resp.Status)
	c.log.Debugf("Response Status Code: %d", resp.StatusCode)
	c.log.Debugf("Response Headers: %s", resp.Header)

	// log message content in DEBUG mode
	// get info of debug mode directly from logger
	if c.log.IsEnabledFor(logging.DEBUG) {
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body = ioutil.NopCloser(bytes.NewReader(body))
		c.log.Debugf("Response Body: %s", body)
	}

	return resp, nil
}

func (c *Client) successfulResponse(resp *http.Response) (*http.Response, error) {
	if resp.StatusCode < 200 || resp.StatusCode > 201 {
		return resp, &ApiError{Response: resp}
	}
	return resp, nil
}

func (c *Client) get(path string) (*http.Response, error) {
	return c.execute("GET", path, nil)
}

func (c *Client) getSuccessful(path string) (*http.Response, error) {
	resp, err := c.execute("GET", path, nil)
	if err != nil {
		return resp, err
	}
	return c.successfulResponse(resp)
}

func (c *Client) postSuccessful(path string, body *[]byte) (*http.Response, error) {
	resp, err := c.execute("POST", path, body)
	if err != nil {
		return resp, err
	}
	return c.successfulResponse(resp)
}

func (c *Client) deleteSuccessful(path string) (*http.Response, error) {
	resp, err := c.execute("DELETE", path, nil)
	if err != nil {
		return resp, err
	}
	return c.successfulResponse(resp)
}

func (c *Client) Login() (error) {

	c.log.Infof("Logging as: %s", c.user + c.password)

	var loginRequest struct {
		User string `json:"email"`
		Password string `json:"password"`
	}

	loginRequest.User = c.user
	loginRequest.Password = c.password

	c.log.Debugf("login request: %v", loginRequest)

	bodyBytes, err := json.Marshal(loginRequest)
	if err != nil {
		return err
	}

	resp, err := c.postSuccessful("login", &bodyBytes)
	if err != nil {
		return err
	}

	var result struct {
		Token string `json:"token"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		c.log.Errorf("Failed to parse json from piot server response. Run " +
					 "client in verbose mode to see raw server communication.")
		return err
	}

	c.token = result.Token

	return nil
}

func (c *Client) GetThings(org *string) ([]Thing, error) {

	var result []Thing

	jsonData := map[string]string{
		"query": `
			{
				things (all: false) {
					id, name, type, alias, enabled, last_seen, store_influxdb, store_mysqldb
				}
			}
		`,
	}
    bodyBytes, err := json.Marshal(jsonData)
    if err != nil {
		return result, err
    }

	resp, err := c.postSuccessful("query", &bodyBytes)
	if err != nil {
		return result, err
	}

	var data struct {
		Data struct {
			Things []Thing  `json:"things"`
		}
	}

	// var gqlResponse interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		c.log.Error(err)
	}

	return data.Data.Things, nil
}

func (c *Client) GetOrgs() ([]Org, error) {

	var result []Org

	jsonData := map[string]string{
		"query": `
			{
				orgs {
					id, name
				}
			}
		`,
	}
    bodyBytes, err := json.Marshal(jsonData)
    if err != nil {
		return result, err
    }

	resp, err := c.postSuccessful("query", &bodyBytes)
	if err != nil {
		return result, err
	}

	var data struct {
		Data struct {
			Orgs []Org`json:"orgs"`
		}
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		c.log.Error(err)
	}

	return data.Data.Orgs, nil
}
