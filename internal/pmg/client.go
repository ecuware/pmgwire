package pmg

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	BaseURL    string
	AuthToken  string
	CSRFToken  string
	HTTPClient *http.Client
	Insecure   bool
}

type AuthConfig struct {
	Host     string `yaml:"host"`
	Token    string `yaml:"token"`
	Insecure bool   `yaml:"insecure"`
}

func DefaultAuthConfig() AuthConfig {
	return AuthConfig{
		Host:     "https://localhost:8006",
		Token:    "",
		Insecure: false,
	}
}

func NewClient(cfg AuthConfig) (*Client, error) {
	baseURL := cfg.Host
	if !strings.HasPrefix(baseURL, "http") {
		baseURL = "https://" + baseURL
	}
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid host URL: %w", err)
	}
	if u.Port() == "" {
		u.Host = u.Host + ":8006"
	}
	baseURL = u.String()

	transport := &http.Transport{}
	if cfg.Insecure {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	client := &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Transport: transport,
			Timeout:   30 * time.Second,
		},
		Insecure: cfg.Insecure,
	}

	if cfg.Token != "" {
		client.AuthToken = cfg.Token
	}

	return client, nil
}

func (c *Client) AuthenticateWithTicket(username, password string) error {
	creds := map[string]string{
		"username": username,
		"password": password,
	}
	body, _ := json.Marshal(creds)

	resp, err := c.HTTPClient.Post(c.BaseURL+"/api2/json/access/ticket", "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to request ticket: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("auth failed (status %d): %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Data struct {
			Ticket    string `json:"ticket"`
			CSRFToken string `json:"CSRFPreventionToken"`
			Username  string `json:"username"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode auth response: %w", err)
	}

	c.AuthToken = result.Data.Ticket
	c.CSRFToken = result.Data.CSRFToken
	return nil
}

func (c *Client) doRequest(method, path string, body interface{}) (json.RawMessage, error) {
	var reqBody io.Reader
	if body != nil {
		formValues := url.Values{}
		switch v := body.(type) {
		case map[string]string:
			for key, val := range v {
				formValues.Set(key, val)
			}
		case map[string]interface{}:
			for key, val := range v {
				formValues.Set(key, fmt.Sprintf("%v", val))
			}
		}
		reqBody = strings.NewReader(formValues.Encode())
	}

	req, err := http.NewRequest(method, c.BaseURL+"/api2/json"+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	if c.AuthToken != "" {
		req.AddCookie(&http.Cookie{Name: "PVEAuthCookie", Value: c.AuthToken})
	}
	if c.CSRFToken != "" && (method == "POST" || method == "PUT" || method == "DELETE") {
		req.Header.Set("CSRFPreventionToken", c.CSRFToken)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return respBody, nil
	}

	return result.Data, nil
}

func (c *Client) Get(path string) (json.RawMessage, error) {
	return c.doRequest("GET", path, nil)
}

func (c *Client) Post(path string, body interface{}) (json.RawMessage, error) {
	return c.doRequest("POST", path, body)
}

func (c *Client) Put(path string, body interface{}) (json.RawMessage, error) {
	return c.doRequest("PUT", path, body)
}

func (c *Client) Delete(path string) (json.RawMessage, error) {
	return c.doRequest("DELETE", path, nil)
}