package pmg

import (
	"encoding/json"
	"fmt"
)

type WhoEmail struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type WhoDomain struct {
	ID     int    `json:"id"`
	Domain string `json:"domain"`
}

func (c *Client) ListWhoEmails(whoID int) ([]WhoEmail, error) {
	path := fmt.Sprintf("/config/ruledb/who/%d/email", whoID)
	data, err := c.Get(path)
	if err != nil {
		return nil, fmt.Errorf("listing who emails for who_id %d: %w", whoID, err)
	}

	var emails []WhoEmail
	if err := json.Unmarshal(data, &emails); err != nil {
		return nil, fmt.Errorf("parsing who email response: %w", err)
	}

	return emails, nil
}

func (c *Client) ListWhoDomains(whoID int) ([]WhoDomain, error) {
	path := fmt.Sprintf("/config/ruledb/who/%d/domain", whoID)
	data, err := c.Get(path)
	if err != nil {
		return nil, fmt.Errorf("listing who domains for who_id %d: %w", whoID, err)
	}

	var domains []WhoDomain
	if err := json.Unmarshal(data, &domains); err != nil {
		return nil, fmt.Errorf("parsing who domain response: %w", err)
	}

	return domains, nil
}

func (c *Client) AddWhoEmail(whoID int, email string) error {
	path := fmt.Sprintf("/config/ruledb/who/%d/email", whoID)
	_, err := c.Post(path, map[string]interface{}{
		"email": email,
	})
	if err != nil {
		return fmt.Errorf("adding email %s to who_id %d: %w", email, whoID, err)
	}
	return nil
}

func (c *Client) AddWhoDomain(whoID int, domain string) error {
	path := fmt.Sprintf("/config/ruledb/who/%d/domain", whoID)
	_, err := c.Post(path, map[string]interface{}{
		"domain": domain,
	})
	if err != nil {
		return fmt.Errorf("adding domain %s to who_id %d: %w", domain, whoID, err)
	}
	return nil
}

func (c *Client) RemoveWhoEmail(whoID int, emailID int) error {
	path := fmt.Sprintf("/config/ruledb/who/%d/email/%d", whoID, emailID)
	_, err := c.Delete(path)
	if err != nil {
		return fmt.Errorf("removing email id %d from who_id %d: %w", emailID, whoID, err)
	}
	return nil
}

func (c *Client) RemoveWhoDomain(whoID int, domainID int) error {
	path := fmt.Sprintf("/config/ruledb/who/%d/domain/%d", whoID, domainID)
	_, err := c.Delete(path)
	if err != nil {
		return fmt.Errorf("removing domain id %d from who_id %d: %w", domainID, whoID, err)
	}
	return nil
}

func (c *Client) ListWhoObjects() (json.RawMessage, error) {
	return c.Get("/config/ruledb/who")
}