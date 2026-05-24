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
	Register(&RuleDBWhatListAction{})
	Register(&RuleDBWhatAddAction{})
	Register(&RuleDBWhatRemoveAction{})
	Register(&RuleDBRuleListAction{})
	Register(&RuleDBRuleCreateAction{})
	Register(&RuleDBRuleRemoveAction{})
}

type RuleDBWhatListAction struct{}

func (a *RuleDBWhatListAction) Name() string { return "ruledb.what.list" }

func (a *RuleDBWhatListAction) Execute(ctx context.Context, client *pmg.Client, input Data, params Params, filters map[string]string) (Data, error) {
	whatID := getIntParam(params, "what_id")

	if whatID > 0 {
		contentTypes, err := client.ListWhatContentTypes(whatID)
		if err != nil {
			return nil, err
		}
		patterns, err := client.ListWhatPatterns(whatID)
		if err != nil {
			return nil, err
		}
		fields, err := client.ListWhatFields(whatID)
		if err != nil {
			return nil, err
		}

		ctList := make([]string, 0, len(contentTypes))
		for _, ct := range contentTypes {
			ctList = append(ctList, ct.ContentType)
		}
		patList := make([]string, 0, len(patterns))
		for _, p := range patterns {
			patList = append(patList, p.Pattern)
		}
		fieldList := make([]string, 0, len(fields))
		for _, f := range fields {
			fieldList = append(fieldList, f.Field)
		}

		fmt.Printf("  What Object %d: %d content types, %d patterns, %d fields\n",
			whatID, len(contentTypes), len(patterns), len(fields))

		return Data{
			"what_id":        whatID,
			"content_types":   ctList,
			"patterns":        patList,
			"fields":          fieldList,
		}, nil
	}

	whats, err := client.ListWhatObjects()
	if err != nil {
		return nil, err
	}

	fmt.Printf("  Found %d What Objects\n", len(whats))

	whatList := make([]map[string]interface{}, 0, len(whats))
	for _, w := range whats {
		whatList = append(whatList, map[string]interface{}{
			"id":   w.ID,
			"name": w.Name,
		})
	}

	return Data{
		"what_objects": whatList,
		"count":        len(whats),
	}, nil
}

type RuleDBWhatAddAction struct{}

func (a *RuleDBWhatAddAction) Name() string { return "ruledb.what.add" }

func (a *RuleDBWhatAddAction) Execute(ctx context.Context, client *pmg.Client, input Data, params Params, filters map[string]string) (Data, error) {
	whatID := getIntParam(params, "what_id")
	mode := getStringParam(params, "mode")
	create := getStringParam(params, "create") == "true"
	name := getStringParam(params, "name")

	if create && whatID == 0 {
		result, err := client.CreateWhatObject(name)
		if err != nil {
			return nil, fmt.Errorf("creating what object: %w", err)
		}
		fmt.Printf("  Created What Object: %s\n", name)
		return Data{
			"name":    name,
			"created": true,
			"raw":     string(result),
		}, nil
	}

	if whatID == 0 {
		return nil, fmt.Errorf("what_id is required when not creating a new object")
	}

	entries := getEntriesFromInput(input)
	if len(entries) == 0 && mode != "create" {
		return nil, fmt.Errorf("no entries provided to add")
	}

	success := 0
	failed := 0
	var addedContentTypes []string
	var addedPatterns []string
	var addedFields []string

	for _, entry := range entries {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}

		var err error
		switch mode {
		case "contenttype":
			err = client.AddWhatContentType(whatID, entry)
			if err == nil {
				addedContentTypes = append(addedContentTypes, entry)
				fmt.Printf("  [OK] Content type added: %s\n", entry)
			}
		case "pattern":
			err = client.AddWhatPattern(whatID, entry)
			if err == nil {
				addedPatterns = append(addedPatterns, entry)
				fmt.Printf("  [OK] Pattern added: %s\n", entry)
			}
		case "field":
			err = client.AddWhatField(whatID, entry)
			if err == nil {
				addedFields = append(addedFields, entry)
				fmt.Printf("  [OK] Field added: %s\n", entry)
			}
		default:
			if strings.Contains(entry, "/") {
				err = client.AddWhatContentType(whatID, entry)
				if err == nil {
					addedContentTypes = append(addedContentTypes, entry)
					fmt.Printf("  [OK] Content type added (auto): %s\n", entry)
				}
			} else if strings.HasPrefix(entry, "X-") || !strings.Contains(entry, ".") {
				err = client.AddWhatField(whatID, entry)
				if err == nil {
					addedFields = append(addedFields, entry)
					fmt.Printf("  [OK] Field added (auto): %s\n", entry)
				}
			} else {
				err = client.AddWhatPattern(whatID, entry)
				if err == nil {
					addedPatterns = append(addedPatterns, entry)
					fmt.Printf("  [OK] Pattern added (auto): %s\n", entry)
				}
			}
		}

		if err != nil {
			fmt.Printf("  [FAIL] %s: %v\n", entry, err)
			failed++
		} else {
			success++
		}
	}

	fmt.Printf("  Summary: %d added, %d failed\n", success, failed)

	return Data{
		"what_id":              whatID,
		"added_content_types":   addedContentTypes,
		"added_patterns":        addedPatterns,
		"added_fields":         addedFields,
		"success":              success,
		"failed":               failed,
	}, nil
}

type RuleDBWhatRemoveAction struct{}

func (a *RuleDBWhatRemoveAction) Name() string { return "ruledb.what.remove" }

func (a *RuleDBWhatRemoveAction) Execute(ctx context.Context, client *pmg.Client, input Data, params Params, filters map[string]string) (Data, error) {
	whatID := getIntParam(params, "what_id")
	deleteObject := getStringParam(params, "delete") == "true"

	if deleteObject {
		if err := client.DeleteWhatObject(whatID); err != nil {
			return nil, err
		}
		fmt.Printf("  [OK] Deleted What Object %d\n", whatID)
		return Data{"what_id": whatID, "deleted": true}, nil
	}

	mode := getStringParam(params, "mode")
	entries := getEntriesFromInput(input)

	if len(entries) == 0 {
		return nil, fmt.Errorf("no entries provided; set delete=true to remove the entire object")
	}

	success := 0
	failed := 0

	switch mode {
	case "contenttype":
		allCTs, err := client.ListWhatContentTypes(whatID)
		if err != nil {
			return nil, err
		}
		for _, entry := range entries {
			for _, ct := range allCTs {
				if ct.ContentType == entry {
					if err := client.RemoveWhatContentType(whatID, ct.ID); err != nil {
						fmt.Printf("  [FAIL] Could not remove content type %s: %v\n", entry, err)
						failed++
					} else {
						fmt.Printf("  [OK] Removed content type: %s\n", entry)
						success++
					}
				}
			}
		}
	case "pattern":
		allPats, err := client.ListWhatPatterns(whatID)
		if err != nil {
			return nil, err
		}
		for _, entry := range entries {
			for _, p := range allPats {
				if p.Pattern == entry {
					if err := client.RemoveWhatPattern(whatID, p.ID); err != nil {
						fmt.Printf("  [FAIL] Could not remove pattern %s: %v\n", entry, err)
						failed++
					} else {
						fmt.Printf("  [OK] Removed pattern: %s\n", entry)
						success++
					}
				}
			}
		}
	case "field":
		allFields, err := client.ListWhatFields(whatID)
		if err != nil {
			return nil, err
		}
		for _, entry := range entries {
			for _, f := range allFields {
				if f.Field == entry {
					if err := client.RemoveWhatField(whatID, f.ID); err != nil {
						fmt.Printf("  [FAIL] Could not remove field %s: %v\n", entry, err)
						failed++
					} else {
						fmt.Printf("  [OK] Removed field: %s\n", entry)
						success++
					}
				}
			}
		}
	default:
		return nil, fmt.Errorf("mode is required (contenttype, pattern, or field)")
	}

	fmt.Printf("  Summary: %d removed, %d failed\n", success, failed)

	return Data{
		"what_id": whatID,
		"success": success,
		"failed":  failed,
	}, nil
}

type RuleDBRuleListAction struct{}

func (a *RuleDBRuleListAction) Name() string { return "ruledb.rule.list" }

func (a *RuleDBRuleListAction) Execute(ctx context.Context, client *pmg.Client, input Data, params Params, filters map[string]string) (Data, error) {
	ruleID := getIntParam(params, "rule_id")

	if ruleID > 0 {
		data, err := client.GetRule(ruleID)
		if err != nil {
			return nil, err
		}
		fmt.Printf("  Rule %d details\n", ruleID)
		return Data{
			"rule_id": ruleID,
			"raw":      string(data),
		}, nil
	}

	rules, err := client.ListRules()
	if err != nil {
		return nil, err
	}

	fmt.Printf("  Found %d rules\n", len(rules))

	ruleList := make([]map[string]interface{}, 0, len(rules))
	for _, r := range rules {
		fmt.Printf("  Rule %d: priority=%d direction=%s action=%s who=%d what=%d\n",
			r.ID, r.Priority, r.Direction, r.Action, r.WhoID, r.WhatID)
		ruleList = append(ruleList, map[string]interface{}{
			"id":        r.ID,
			"priority":  r.Priority,
			"direction": r.Direction,
			"action":    r.Action,
			"who_id":    r.WhoID,
			"what_id":   r.WhatID,
		})
	}

	return Data{
		"rules": ruleList,
		"count": len(rules),
	}, nil
}

type RuleDBRuleCreateAction struct{}

func (a *RuleDBRuleCreateAction) Name() string { return "ruledb.rule.create" }

func (a *RuleDBRuleCreateAction) Execute(ctx context.Context, client *pmg.Client, input Data, params Params, filters map[string]string) (Data, error) {
	ruleParams := make(map[string]interface{})

	priority := getIntParam(params, "priority")
	if priority > 0 {
		ruleParams["priority"] = priority
	}

	direction := getStringParam(params, "direction")
	if direction != "" {
		ruleParams["direction"] = direction
	}

	action := getStringParam(params, "action")
	if action != "" {
		ruleParams["action"] = action
	}

	whoID := getIntParam(params, "who_id")
	if whoID > 0 {
		ruleParams["who_id"] = whoID
	}

	whatID := getIntParam(params, "what_id")
	if whatID > 0 {
		ruleParams["what_id"] = whatID
	}

	name := getStringParam(params, "name")
	if name != "" {
		ruleParams["name"] = name
	}

	result, err := client.CreateRule(ruleParams)
	if err != nil {
		return nil, fmt.Errorf("creating rule: %w", err)
	}

	fmt.Printf("  [OK] Rule created: priority=%d direction=%s action=%s who=%d what=%d\n",
		priority, direction, action, whoID, whatID)

	return Data{
		"priority":  priority,
		"direction": direction,
		"action":    action,
		"who_id":    whoID,
		"what_id":   whatID,
		"raw":       string(result),
	}, nil
}

type RuleDBRuleRemoveAction struct{}

func (a *RuleDBRuleRemoveAction) Name() string { return "ruledb.rule.remove" }

func (a *RuleDBRuleRemoveAction) Execute(ctx context.Context, client *pmg.Client, input Data, params Params, filters map[string]string) (Data, error) {
	ruleID := getIntParam(params, "rule_id")
	if ruleID == 0 {
		return nil, fmt.Errorf("rule_id is required")
	}

	if err := client.DeleteRule(ruleID); err != nil {
		return nil, err
	}

	fmt.Printf("  [OK] Deleted rule %d\n", ruleID)
	return Data{"rule_id": ruleID}, nil
}