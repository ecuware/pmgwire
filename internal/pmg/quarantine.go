package pmg

import (
	"encoding/json"
	"fmt"
	"strings"
)

type QuarantineMail struct {
	ID              string `json:"id"`
	Receiver        string `json:"receiver"`
	From            string `json:"from"`
	EnvelopeSender  string `json:"envelope_sender"`
	Subject         string `json:"subject"`
	SpamScore        float64 `json:"spam_score"`
	Time            string `json:"time"`
}

func (c *Client) ListQuarantineSpam() ([]QuarantineMail, error) {
	data, err := c.Get("/quarantine/spam")
	if err != nil {
		return nil, fmt.Errorf("listing spam quarantine: %w", err)
	}

	var mails []QuarantineMail
	if err := json.Unmarshal(data, &mails); err != nil {
		return nil, fmt.Errorf("parsing spam quarantine response: %w", err)
	}

	return mails, nil
}

func (c *Client) ListQuarantineVirus() ([]QuarantineMail, error) {
	data, err := c.Get("/quarantine/virus")
	if err != nil {
		return nil, fmt.Errorf("listing virus quarantine: %w", err)
	}

	var mails []QuarantineMail
	if err := json.Unmarshal(data, &mails); err != nil {
		return nil, fmt.Errorf("parsing virus quarantine response: %w", err)
	}

	return mails, nil
}

func (c *Client) DeliverQuarantine(id string) error {
	_, err := c.Post("/quarantine/content", map[string]interface{}{
		"id":     id,
		"action": "deliver",
	})
	if err != nil {
		return fmt.Errorf("delivering quarantine mail %s: %w", id, err)
	}
	return nil
}

func (c *Client) DeleteQuarantine(id string) error {
	_, err := c.Post("/quarantine/content", map[string]interface{}{
		"id":     id,
		"action": "delete",
	})
	if err != nil {
		return fmt.Errorf("deleting quarantine mail %s: %w", id, err)
	}
	return nil
}

func FilterQuarantine(mails []QuarantineMail, senderPattern, receiverPattern string) []QuarantineMail {
	var filtered []QuarantineMail
	for _, m := range mails {
		if senderPattern != "" && senderPattern != "*" {
			if !matchPattern(m.From, senderPattern) && !matchPattern(m.EnvelopeSender, senderPattern) {
				continue
			}
		}
		if receiverPattern != "" && receiverPattern != "*" {
			if !matchPattern(m.Receiver, receiverPattern) {
				continue
			}
		}
		filtered = append(filtered, m)
	}
	return filtered
}

func matchPattern(s, pattern string) bool {
	if pattern == "" || pattern == "*" {
		return true
	}
	s = strings.ToLower(s)
	pattern = strings.ToLower(pattern)
	if strings.Contains(pattern, "*") {
		return simpleMatch(s, pattern)
	}
	return strings.Contains(s, pattern)
}

func simpleMatch(s, pattern string) bool {
	pi := 0
	si := 0
	for pi < len(pattern) && si < len(s) {
		if pattern[pi] == '*' {
			pi++
			if pi == len(pattern) {
				return true
			}
			for si < len(s) {
				if simpleMatch(s[si:], pattern[pi:]) {
					return true
				}
				si++
			}
			return false
		}
		if pattern[pi] != s[si] {
			return false
		}
		pi++
		si++
	}
	for pi < len(pattern) && pattern[pi] == '*' {
		pi++
	}
	return pi == len(pattern) && si == len(s)
}