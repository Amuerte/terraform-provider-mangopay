package mangopay

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) GetToken() (*AuthResponse, error) {
	v := url.Values{
		"grant_type": {"client_credentials"},
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/oauth/token", c.Host), strings.NewReader(v.Encode()))
	if err != nil {
		return nil, err
	}

	if c.AuthConfig.ClientID == "" || c.AuthConfig.ClientSecret == "" {
		return nil, fmt.Errorf("define ClientID and ClientSecret in order to get token")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	req.SetBasicAuth(c.AuthConfig.ClientID, c.AuthConfig.ClientSecret)
	req.BasicAuth()

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	ar := AuthResponse{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return nil, err
	}

	return &ar, nil
}
