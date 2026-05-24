package pmg

import (
	"encoding/json"
	"fmt"
)

func (c *Client) GetConfig(section string) (json.RawMessage, error) {
	path := fmt.Sprintf("/config/%s", section)
	data, err := c.Get(path)
	if err != nil {
		return nil, fmt.Errorf("getting config section %s: %w", section, err)
	}
	return data, nil
}

func (c *Client) UpdateConfig(section string, values map[string]interface{}) error {
	path := fmt.Sprintf("/config/%s", section)
	_, err := c.Put(path, values)
	if err != nil {
		return fmt.Errorf("updating config section %s: %w", section, err)
	}
	return nil
}