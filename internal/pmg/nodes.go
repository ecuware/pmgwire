package pmg

import (
	"encoding/json"
	"fmt"
)

func (c *Client) GetNodes() (json.RawMessage, error) {
	data, err := c.Get("/nodes")
	if err != nil {
		return nil, fmt.Errorf("getting nodes: %w", err)
	}
	return data, nil
}

func (c *Client) GetNodeStatus(node string) (json.RawMessage, error) {
	path := fmt.Sprintf("/nodes/%s/status", node)
	data, err := c.Get(path)
	if err != nil {
		return nil, fmt.Errorf("getting status for node %s: %w", node, err)
	}
	return data, nil
}