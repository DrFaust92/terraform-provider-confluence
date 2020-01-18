package confluence

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client provides a connection to the Confluence API
type Client struct {
	client    *http.Client
	baseURL   *url.URL
	publicURL *url.URL
}

// NewClientInput provides information to connect to the Confluence API
type NewClientInput struct {
	site  string
	user  string
	token string
}

// ErrorResponse describes why a request failed
type ErrorResponse struct {
	StatusCode int `json:"statusCode,omitempty"`
	Data       struct {
		Authorized bool     `json:"authorized,omitempty"`
		Valid      bool     `json:"valid,omitempty"`
		Errors     []string `json:"errors,omitempty"`
		Successful bool     `json:"successful,omitempty"`
	} `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

// NewClient returns an authenticated client ready to use
func NewClient(input *NewClientInput) *Client {
	publicURL := url.URL{
		Scheme: "https",
		Host:   input.site + ".atlassian.net",
	}
	baseURL := publicURL
	baseURL.User = url.UserPassword(input.user, input.token)
	return &Client{
		client: &http.Client{
			Timeout: time.Second * 10,
		},
		baseURL:   &baseURL,
		publicURL: &publicURL,
	}
}

// Post uses the client to send a POST request
func (c *Client) Post(path string, body interface{}, result interface{}) error {
	return c.do("POST", path, body, result)
}

// Get uses the client to send a GET request
func (c *Client) Get(path string, result interface{}) error {
	return c.do("GET", path, nil, result)
}

// Put uses the client to send a PUT request
func (c *Client) Put(path string, body interface{}, result interface{}) error {
	return c.do("PUT", path, body, result)
}

// Delete uses the client to send a DELETE request
func (c *Client) Delete(path string) error {
	return c.do("DELETE", path, nil, nil)
}

// do uses the client to send a specified request
func (c *Client) do(method string, path string, body interface{}, result interface{}) error {
	u, err := c.baseURL.Parse(path)
	if err != nil {
		return err
	}
	var bodyReader io.Reader
	var bodyBytes []byte
	if body != nil {
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}
	req, err := http.NewRequest(method, u.String(), bodyReader)
	if err != nil {
		return err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	expectedStatusCode := map[string]int{
		"POST":   200,
		"PUT":    200,
		"GET":    200,
		"DELETE": 204,
	}
	if resp.StatusCode != expectedStatusCode[method] {
		var errResponse ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&errResponse)
		if err != nil {
			return fmt.Errorf("%s\n\n%s %s\n%s\n\n%v",
				resp.Status, method, path, string(bodyBytes), err)
		}
		return fmt.Errorf("%s\n\n%s %s\n%s\n\n%s",
			resp.Status, method, path, string(bodyBytes), &errResponse)
	}
	if result != nil {
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *ErrorResponse) String() string {
	d := e.Data
	var errorsString string
	if len(d.Errors) > 0 {
		errorsString = fmt.Sprintf("\n  * %s", strings.Join(d.Errors, "\n  * "))
	}
	return fmt.Sprintf("%s\nAuthorized: %t\nValid: %t\nSuccessful: %t%s",
		e.Message, d.Authorized, d.Valid, d.Successful, errorsString)
}

// URL returns the public URL for a given path
func (c *Client) URL(path string) string {
	u, err := c.publicURL.Parse(path)
	if err != nil {
		return ""
	}
	return u.String()
}
