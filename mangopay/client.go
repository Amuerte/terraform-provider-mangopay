package mangopay

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"time"
)

type Client struct {
	Host                string
	HTTPClient          *http.Client
	AuthConfig          OAuth2Config
	MangopayEnvironment string
	Token               string
}

type OAuth2Config struct {
	ClientID     string
	ClientSecret string
}

func New(clientID, clientSecret, environment *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second, Transport: &loggingTransport{}},
	}

	c.Host = "https://api.sandbox.mangopay.com/v2.01"

	// If username or password not provided, return empty client
	if clientID == nil || clientSecret == nil {
		return &c, nil
	}

	c.AuthConfig = OAuth2Config{
		ClientID:     *clientID,
		ClientSecret: *clientSecret,
	}

	ar, err := c.GetToken()
	if err != nil {
		return nil, err
	}

	c.Token = "Bearer " + ar.AccessToken

	return &c, nil
}

type Response struct {
	Message string `json:"message"`
}

func (c *Client) doRequest(req *http.Request, authToken *string) ([]byte, error) {
	token := c.Token

	if authToken != nil {
		token = *authToken
	}

	if req.Header["Authorization"] == nil {
		req.Header.Set("Authorization", token)
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}

func (s *loggingTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	bytes, _ := httputil.DumpRequestOut(r, true)

	resp, err := http.DefaultTransport.RoundTrip(r)
	// err is returned after dumping the response

	respBytes, _ := httputil.DumpResponse(resp, true)
	bytes = append(bytes, respBytes...)

	fmt.Printf("%s\n", bytes)

	return resp, err
}

type loggingTransport struct{}
