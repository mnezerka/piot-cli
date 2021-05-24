package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/op/go-logging"
	"github.com/spf13/viper"
)

type Client struct {
	user      string
	password  string
	piot_host string
	log       *logging.Logger
	token     string
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

	c.log.Debugf("------%s Request to: %s", method, url)
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

	if c.token == "" {
		c.log.Debug("Setting basic authorization (header)")
		req.SetBasicAuth(c.user, c.password)
	} else {
		c.log.Debugf("Setting bearer authorization (reusing token %s)", c.token)
		req.Header.Add("Authorization", "Bearer "+c.token)
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

func (c *Client) gqlQuerySuccessful(gql string) (*http.Response, error) {

	jsonData := map[string]string{
		"query": fmt.Sprintf(`%s`, gql),
	}

	bodyBytes, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}

	resp, err := c.postSuccessful("query", &bodyBytes)
	if err != nil {
		return nil, err
	}

	// check gql error inside successfull response
	var data struct {
		Errors []GqlError `json:"errors"`
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	if len(data.Errors) > 0 {

		var errorMessage string

		for i := 0; i < len(data.Errors); i++ {

			if i > 0 {
				errorMessage += "\n"
			}

			errorMessage += fmt.Sprintf("GraphQL error: msg='%s'", data.Errors[i].Message)
			for j := 0; j < len(data.Errors[i].Locations); j++ {
				errorMessage += fmt.Sprintf(" line=%d, col=%d",
					data.Errors[i].Locations[j].Line,
					data.Errors[i].Locations[j].Column,
				)
			}
		}

		return resp, fmt.Errorf("%s", errorMessage)
	}

	resp.Body = ioutil.NopCloser(bytes.NewReader(body))

	return resp, nil
}

func (c *Client) Login() error {

	if c.token != "" {
		c.log.Debug("Reusing existing token")
	}

	c.log.Infof("Logging as: %s", c.user)

	var loginRequest struct {
		User     string `json:"email"`
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

func (c *Client) GetThings(all bool) ([]Thing, error) {

	var result []Thing

	all_str := "false"
	if all {
		all_str = "true"
	}

	gql := fmt.Sprintf(`
		{
			things (all: %s) {
				id, name, type, alias, enabled, last_seen, last_seen_interval, store_influxdb, store_mysqldb
			}
		}
		`, all_str)
	resp, err := c.gqlQuerySuccessful(gql)
	if err != nil {
		return result, err
	}

	var data struct {
		Data struct {
			Things []Thing `json:"things"`
		}
	}

	// var gqlResponse interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		c.log.Error(err)
	}

	return data.Data.Things, nil
}

type OrgFilterFunctionType = func(s *Org) bool

func (c *Client) GetOrgs(filter OrgFilterFunctionType) ([]Org, error) {

	var result []Org

	gql := "{orgs {id, name}}"

	resp, err := c.gqlQuerySuccessful(gql)
	if err != nil {
		return result, err
	}

	var data struct {
		Data struct {
			Orgs []Org `json:"orgs"`
		}
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return result, err
	}

	result = data.Data.Orgs

	// apply filtering if filter function was provided
	if filter != nil {
		var filteredResult []Org
		for _, o := range data.Data.Orgs {
			if filter(&o) {
				filteredResult = append(filteredResult, o)
			}
		}
		result = filteredResult
	}

	return result, nil
}

func (c *Client) GetOrgByName(name string) (*Org, error) {

	orgs, err := c.GetOrgs(func(org *Org) bool { return org.Name == name })

	if err != nil {
		return nil, err
	}

	if len(orgs) > 0 {
		return &orgs[0], nil
	}

	return nil, nil
}

func (c *Client) SetCurrentOrg(name string) error {

	org, err := c.GetOrgByName(name)
	if err != nil {
		return err
	}

	if org == nil {
		return fmt.Errorf("Organization '%s' does not exist", name)
	}

	gql := fmt.Sprintf(`mutation { updateUserProfile(profile: {org_id: "%s"}) {org_id}}`, org.Id)
	resp, err := c.gqlQuerySuccessful(gql)

	c.log.Infof("result of set: %v", resp.Body)
	return err
}

func (c *Client) CreateThing(name, thing_type string) (string, error) {

	c.log.Infof("Creating new thing: name='%s', type='%s'", name, thing_type)

	gql := fmt.Sprintf(`
		mutation {
			createThing(name: "%s", type: "%s") {id}
		}
	`, name, thing_type)

	resp, err := c.gqlQuerySuccessful(gql)
	if err != nil {
		return "", err
	}

	c.log.Infof("%v", resp)

	/*
		var data struct {
			Data struct {
				Thing Thing `json:"thing"`
			}
		}

		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return "", err
		}

		return data.Thing.Id, nil
	*/
	return "id", nil
}

func (c *Client) DeleteThing() error {

	jsonData := map[string]string{
		"query": fmt.Sprintf(`
			mutation {
				deleteThing(id : "%s")
			}
		`, "idtodel"),
	}

	c.log.Infof("%v", jsonData)

	return nil
}

func (c *Client) GetUserProfile() (UserProfile, error) {

	var result UserProfile

	gql := `
		{
			userProfile {
				email, is_admin, org_id, orgs {id, name} 
			}
		}
	`
	resp, err := c.gqlQuerySuccessful(gql)
	if err != nil {
		return result, err
	}

	var data struct {
		Data struct {
			Profile UserProfile `json:"userProfile"`
		}
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		c.log.Error(err)
	}

	return data.Data.Profile, nil
}
