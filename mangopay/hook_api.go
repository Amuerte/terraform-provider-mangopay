package mangopay

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetAllHooks - Returns a list of configured hooks
func (c *Client) GetAllHooks() ([]Hook, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/hooks", c.Host, c.AuthConfig.ClientID), nil)

	hooks := []Hook{}
	if err != nil {
		return hooks, err
	}

	body, err := c.doRequest(req, nil)
	if err != nil {
		return hooks, err
	}

	err = json.Unmarshal(body, &hooks)
	if err != nil {
		return []Hook{}, err
	}

	return hooks, nil
}

// GetHooks - Returns a list of configured hooks
func (c *Client) GetHook(hookID int64) (Hook, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/hooks/%d", c.Host, c.AuthConfig.ClientID, hookID), nil)

	hook := Hook{}
	if err != nil {
		return hook, err
	}

	body, err := c.doRequest(req, nil)
	if err != nil {
		return hook, err
	}

	err = json.Unmarshal(body, &hook)
	if err != nil {
		return Hook{}, err
	}

	return hook, nil
}
