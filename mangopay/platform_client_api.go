package mangopay

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetPlatformClient - Returns a client detail
func (c *Client) GetPlatformClient() (PlatformClient, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/clients", c.Host, c.AuthConfig.ClientID), nil)

	platformClient := PlatformClient{}
	if err != nil {
		return platformClient, err
	}

	body, err := c.doRequest(req, nil)
	if err != nil {
		return platformClient, err
	}

	err = json.Unmarshal(body, &platformClient)
	if err != nil {
		return PlatformClient{}, err
	}

	return platformClient, nil
}
