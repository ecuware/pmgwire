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

type WhatContentType struct {
	ID          int    `json:"id"`
	ContentType string `json:"contenttype"`
}

type WhatPattern struct {
	ID      int    `json:"id"`
	Pattern string `json:"pattern"`
}

type WhatField struct {
	ID    int    `json:"id"`
	Field string `json:"field"`
}

type WhatObject struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (c *Client) ListWhatObjects() ([]WhatObject, error) {
	data, err := c.Get("/config/ruledb/what")
	if err != nil {
		return nil, fmt.Errorf("listing what objects: %w", err)
	}

	var whats []WhatObject
	if err := json.Unmarshal(data, &whats); err != nil {
		return nil, fmt.Errorf("parsing what objects response: %w", err)
	}

	return whats, nil
}

func (c *Client) GetWhatObject(whatID int) (json.RawMessage, error) {
	path := fmt.Sprintf("/config/ruledb/what/%d", whatID)
	data, err := c.Get(path)
	if err != nil {
		return nil, fmt.Errorf("getting what object %d: %w", whatID, err)
	}
	return data, nil
}

func (c *Client) CreateWhatObject(name string) (json.RawMessage, error) {
	data, err := c.Post("/config/ruledb/what", map[string]interface{}{
		"name": name,
	})
	if err != nil {
		return nil, fmt.Errorf("creating what object %s: %w", name, err)
	}
	return data, nil
}

func (c *Client) DeleteWhatObject(whatID int) error {
	path := fmt.Sprintf("/config/ruledb/what/%d", whatID)
	_, err := c.Delete(path)
	if err != nil {
		return fmt.Errorf("deleting what object %d: %w", whatID, err)
	}
	return nil
}

func (c *Client) ListWhatContentTypes(whatID int) ([]WhatContentType, error) {
	path := fmt.Sprintf("/config/ruledb/what/%d/contenttype", whatID)
	data, err := c.Get(path)
	if err != nil {
		return nil, fmt.Errorf("listing content types for what_id %d: %w", whatID, err)
	}

	var cts []WhatContentType
	if err := json.Unmarshal(data, &cts); err != nil {
		return nil, fmt.Errorf("parsing content types response: %w", err)
	}
	return cts, nil
}

func (c *Client) AddWhatContentType(whatID int, contentType string) error {
	path := fmt.Sprintf("/config/ruledb/what/%d/contenttype", whatID)
	_, err := c.Post(path, map[string]interface{}{
		"contenttype": contentType,
	})
	if err != nil {
		return fmt.Errorf("adding content type %s to what_id %d: %w", contentType, whatID, err)
	}
	return nil
}

func (c *Client) RemoveWhatContentType(whatID int, ctID int) error {
	path := fmt.Sprintf("/config/ruledb/what/%d/contenttype/%d", whatID, ctID)
	_, err := c.Delete(path)
	if err != nil {
		return fmt.Errorf("removing content type %d from what_id %d: %w", ctID, whatID, err)
	}
	return nil
}

func (c *Client) ListWhatPatterns(whatID int) ([]WhatPattern, error) {
	path := fmt.Sprintf("/config/ruledb/what/%d/pattern", whatID)
	data, err := c.Get(path)
	if err != nil {
		return nil, fmt.Errorf("listing patterns for what_id %d: %w", whatID, err)
	}

	var patterns []WhatPattern
	if err := json.Unmarshal(data, &patterns); err != nil {
		return nil, fmt.Errorf("parsing patterns response: %w", err)
	}
	return patterns, nil
}

func (c *Client) AddWhatPattern(whatID int, pattern string) error {
	path := fmt.Sprintf("/config/ruledb/what/%d/pattern", whatID)
	_, err := c.Post(path, map[string]interface{}{
		"pattern": pattern,
	})
	if err != nil {
		return fmt.Errorf("adding pattern %s to what_id %d: %w", pattern, whatID, err)
	}
	return nil
}

func (c *Client) RemoveWhatPattern(whatID int, patternID int) error {
	path := fmt.Sprintf("/config/ruledb/what/%d/pattern/%d", whatID, patternID)
	_, err := c.Delete(path)
	if err != nil {
		return fmt.Errorf("removing pattern %d from what_id %d: %w", patternID, whatID, err)
	}
	return nil
}

func (c *Client) ListWhatFields(whatID int) ([]WhatField, error) {
	path := fmt.Sprintf("/config/ruledb/what/%d/field", whatID)
	data, err := c.Get(path)
	if err != nil {
		return nil, fmt.Errorf("listing fields for what_id %d: %w", whatID, err)
	}

	var fields []WhatField
	if err := json.Unmarshal(data, &fields); err != nil {
		return nil, fmt.Errorf("parsing fields response: %w", err)
	}
	return fields, nil
}

func (c *Client) AddWhatField(whatID int, field string) error {
	path := fmt.Sprintf("/config/ruledb/what/%d/field", whatID)
	_, err := c.Post(path, map[string]interface{}{
		"field": field,
	})
	if err != nil {
		return fmt.Errorf("adding field %s to what_id %d: %w", field, whatID, err)
	}
	return nil
}

func (c *Client) RemoveWhatField(whatID int, fieldID int) error {
	path := fmt.Sprintf("/config/ruledb/what/%d/field/%d", whatID, fieldID)
	_, err := c.Delete(path)
	if err != nil {
		return fmt.Errorf("removing field %d from what_id %d: %w", fieldID, whatID, err)
	}
	return nil
}

type Rule struct {
	ID        int    `json:"id"`
	Priority  int    `json:"priority"`
	Direction string `json:"direction"`
	Action    string `json:"action"`
	WhoID     int    `json:"who_id"`
	WhatID    int    `json:"what_id"`
	Name      string `json:"name,omitempty"`
}

func (c *Client) ListRules() ([]Rule, error) {
	data, err := c.Get("/config/ruledb/rule")
	if err != nil {
		return nil, fmt.Errorf("listing rules: %w", err)
	}

	var rules []Rule
	if err := json.Unmarshal(data, &rules); err != nil {
		return nil, fmt.Errorf("parsing rules response: %w", err)
	}
	return rules, nil
}

func (c *Client) GetRule(ruleID int) (json.RawMessage, error) {
	path := fmt.Sprintf("/config/ruledb/rule/%d", ruleID)
	data, err := c.Get(path)
	if err != nil {
		return nil, fmt.Errorf("getting rule %d: %w", ruleID, err)
	}
	return data, nil
}

func (c *Client) CreateRule(params map[string]interface{}) (json.RawMessage, error) {
	data, err := c.Post("/config/ruledb/rule", params)
	if err != nil {
		return nil, fmt.Errorf("creating rule: %w", err)
	}
	return data, nil
}

func (c *Client) DeleteRule(ruleID int) error {
	path := fmt.Sprintf("/config/ruledb/rule/%d", ruleID)
	_, err := c.Delete(path)
	if err != nil {
		return fmt.Errorf("deleting rule %d: %w", ruleID, err)
	}
	return nil
}