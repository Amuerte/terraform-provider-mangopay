package mangopay

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// GetAllHooks - Returns a list of configured hooks
func (c *Client) GetAllHooks() ([]Hook, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/hooks", c.Host, c.AuthConfig.ClientID), nil)

	hooks := []Hook{}
	if err != nil {
		return hooks, err
	}

	q := req.URL.Query()
	q.Add("per_page", "100")
	req.URL.RawQuery = q.Encode()

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

func (c *Client) GetHook(hookID string) (Hook, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/hooks/%s", c.Host, c.AuthConfig.ClientID, hookID), nil)

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

func (c *Client) CreateHook(hook Hook) (*Hook, error) {
	rb, err := json.Marshal(hook)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/hooks", c.Host, c.AuthConfig.ClientID), strings.NewReader((string(rb))))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	newHook := Hook{}
	err = json.Unmarshal(body, &newHook)
	if err != nil {
		return &newHook, err
	}

	return &newHook, nil
}

func (c *Client) UpdateHook(hookID string, hook Hook) (*Hook, error) {
	rb, err := json.Marshal(hook)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/%s/hooks/%s", c.Host, c.AuthConfig.ClientID, hookID), strings.NewReader((string(rb))))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	updatedHook := Hook{}
	err = json.Unmarshal(body, &updatedHook)
	if err != nil {
		return &updatedHook, err
	}

	return &updatedHook, nil
}
