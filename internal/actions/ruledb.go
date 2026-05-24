package actions

import (
	"context"
	"fmt"
	"strings"

	"github.com/ecuware/pmgwire/internal/pmg"
)

type RuleDBWhoListAction struct{}

func (a *RuleDBWhoListAction) Name() string { return "ruledb.who.list" }

func (a *RuleDBWhoListAction) Execute(ctx context.Context, client *pmg.Client, input Data, params Params, filters map[string]string) (Data, error) {
	whoID := getIntParam(params, "who_id")

	emails, err := client.ListWhoEmails(whoID)
	if err != nil {
		return nil, err
	}

	domains, err := client.ListWhoDomains(whoID)
	if err != nil {
		return nil, err
	}

	fmt.Printf("  Who Object %d: %d emails, %d domains\n", whoID, len(emails), len(domains))

	emailList := make([]string, 0, len(emails))
	for _, e := range emails {
		emailList = append(emailList, e.Email)
	}

	domainList := make([]string, 0, len(domains))
	for _, d := range domains {
		domainList = append(domainList, d.Domain)
	}

	return Data{
		"who_id":  whoID,
		"emails":  emailList,
		"domains": domainList,
	}, nil
}

type RuleDBWhoAddAction struct{}

func (a *RuleDBWhoAddAction) Name() string { return "ruledb.who.add" }

func (a *RuleDBWhoAddAction) Execute(ctx context.Context, client *pmg.Client, input Data, params Params, filters map[string]string) (Data, error) {
	whoID := getIntParam(params, "who_id")
	mode := getStringParam(params, "mode")
	if mode == "" {
		mode = "auto"
	}

	entries := getEntriesFromInput(input)
	if len(entries) == 0 {
		return nil, fmt.Errorf("no entries provided to add")
	}

	success := 0
	failed := 0
	var addedEmails []string
	var addedDomains []string

	for _, entry := range entries {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}

		isEmail := strings.Contains(entry, "@")
		isDomain := !isEmail && (strings.HasPrefix(entry, "@") || strings.Contains(entry, "."))

		if mode == "auto" {
			if isEmail {
				if err := client.AddWhoEmail(whoID, entry); err != nil {
					fmt.Printf("  [FAIL] Email %s: %v\n", entry, err)
					failed++
				} else {
					fmt.Printf("  [OK] Email added: %s\n", entry)
					addedEmails = append(addedEmails, entry)
					success++
				}
			} else if isDomain {
				domain := strings.TrimPrefix(entry, "@")
				if err := client.AddWhoDomain(whoID, domain); err != nil {
					fmt.Printf("  [FAIL] Domain %s: %v\n", domain, err)
					failed++
				} else {
					fmt.Printf("  [OK] Domain added: %s\n", domain)
					addedDomains = append(addedDomains, domain)
					success++
				}
			} else {
				fmt.Printf("  [SKIP] Unrecognized format: %s\n", entry)
				failed++
			}
		} else if mode == "email" && isEmail {
			if err := client.AddWhoEmail(whoID, entry); err != nil {
				fmt.Printf("  [FAIL] Email %s: %v\n", entry, err)
				failed++
			} else {
				fmt.Printf("  [OK] Email added: %s\n", entry)
				addedEmails = append(addedEmails, entry)
				success++
			}
		} else if mode == "domain" {
			domain := strings.TrimPrefix(entry, "@")
			if err := client.AddWhoDomain(whoID, domain); err != nil {
				fmt.Printf("  [FAIL] Domain %s: %v\n", domain, err)
				failed++
			} else {
				fmt.Printf("  [OK] Domain added: %s\n", domain)
				addedDomains = append(addedDomains, domain)
				success++
			}
		}
	}

	fmt.Printf("  Summary: %d added, %d failed\n", success, failed)

	return Data{
		"who_id":       whoID,
		"added_emails": addedEmails,
		"added_domains": addedDomains,
		"success":      success,
		"failed":       failed,
	}, nil
}

type RuleDBWhoRemoveAction struct{}

func (a *RuleDBWhoRemoveAction) Name() string { return "ruledb.who.remove" }

func (a *RuleDBWhoRemoveAction) Execute(ctx context.Context, client *pmg.Client, input Data, params Params, filters map[string]string) (Data, error) {
	whoID := getIntParam(params, "who_id")

	emails, ok := input["emails"].([]string)
	if ok && len(emails) > 0 {
		allEmails, err := client.ListWhoEmails(whoID)
		if err != nil {
			return nil, err
		}

		removed := 0
		for _, toRemove := range emails {
			for _, existing := range allEmails {
				if existing.Email == toRemove {
					if err := client.RemoveWhoEmail(whoID, existing.ID); err != nil {
						fmt.Printf("  [FAIL] Could not remove email %s: %v\n", toRemove, err)
					} else {
						fmt.Printf("  [OK] Removed email: %s\n", toRemove)
						removed++
					}
				}
			}
		}
		fmt.Printf("  Removed %d emails\n", removed)
	}

	return Data{"who_id": whoID}, nil
}

func getEntriesFromInput(input Data) []string {
	if entries, ok := input["entries"]; ok {
		switch v := entries.(type) {
		case []string:
			return v
		case []interface{}:
			result := make([]string, 0, len(v))
			for _, item := range v {
				result = append(result, fmt.Sprintf("%v", item))
			}
			return result
		}
	}

	if emails, ok := input["emails"]; ok {
		switch v := emails.(type) {
		case []string:
			return v
		case []interface{}:
			result := make([]string, 0, len(v))
			for _, item := range v {
				result = append(result, fmt.Sprintf("%v", item))
			}
			return result
		}
	}

	return nil
}

func getIntParam(params Params, key string) int {
	val, ok := params[key]
	if !ok {
		return 0
	}
	switch v := val.(type) {
	case int:
		return v
	case float64:
		return int(v)
	case string:
		var i int
		fmt.Sscanf(v, "%d", &i)
		return i
	}
	return 0
}

func getStringParam(params Params, key string) string {
	val, ok := params[key]
	if !ok {
		return ""
	}
	return fmt.Sprintf("%v", val)
}

func init() {
	Register(&RuleDBWhoListAction{})
	Register(&RuleDBWhoAddAction{})
	Register(&RuleDBWhoRemoveAction{})
}