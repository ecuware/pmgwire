package pmg

import (
	"encoding/json"
	"fmt"
)

func (c *Client) GetMailStatistics(node string) (json.RawMessage, error) {
	path := fmt.Sprintf("/nodes/%s/statistics/mail", node)
	data, err := c.Get(path)
	if err != nil {
		return nil, fmt.Errorf("getting mail statistics for node %s: %w", node, err)
	}
	return data, nil
}

func (c *Client) GetVirusStatistics(node string) (json.RawMessage, error) {
	path := fmt.Sprintf("/nodes/%s/statistics/virus", node)
	data, err := c.Get(path)
	if err != nil {
		return nil, fmt.Errorf("getting virus statistics for node %s: %w", node, err)
	}
	return data, nil
}