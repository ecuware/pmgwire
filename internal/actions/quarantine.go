package actions

import (
	"context"
	"fmt"

	"github.com/ecuware/pmgwire/internal/pmg"
)

type QuarantineListAction struct{}

func (a *QuarantineListAction) Name() string { return "quarantine.list" }

func (a *QuarantineListAction) Execute(ctx context.Context, client *pmg.Client, input Data, params Params, filters map[string]string) (Data, error) {
	qType := "spam"
	if t, ok := params["type"]; ok {
		qType = fmt.Sprintf("%v", t)
	}

	var mails []pmg.QuarantineMail
	var err error

	switch qType {
	case "virus":
		mails, err = client.ListQuarantineVirus()
	default:
		mails, err = client.ListQuarantineSpam()
	}

	if err != nil {
		return nil, err
	}

	senderFilter := filters["sender"]
	receiverFilter := filters["receiver"]
	filtered := pmg.FilterQuarantine(mails, senderFilter, receiverFilter)

	fmt.Printf("  Found %d quarantine mails (filtered from %d)\n", len(filtered), len(mails))

	result := Data{
		"count": len(filtered),
		"mails": filtered,
		"ids":   extractIDs(filtered),
	}

	return result, nil
}

func extractIDs(mails []pmg.QuarantineMail) []string {
	ids := make([]string, 0, len(mails))
	for _, m := range mails {
		ids = append(ids, m.ID)
	}
	return ids
}

type QuarantineDeliverAction struct{}

func (a *QuarantineDeliverAction) Name() string { return "quarantine.deliver" }

func (a *QuarantineDeliverAction) Execute(ctx context.Context, client *pmg.Client, input Data, params Params, filters map[string]string) (Data, error) {
	ids := getIDsFromInput(input)
	if len(ids) == 0 {
		return nil, fmt.Errorf("no mail IDs provided for delivery")
	}

	delivered := 0
	failed := 0

	for _, id := range ids {
		if err := client.DeliverQuarantine(id); err != nil {
			fmt.Printf("  [FAIL] Could not deliver %s: %v\n", id, err)
			failed++
		} else {
			fmt.Printf("  [OK] Delivered: %s\n", id)
			delivered++
		}
	}

	return Data{
		"delivered": delivered,
		"failed":     failed,
		"total":      len(ids),
	}, nil
}

type QuarantineDeleteAction struct{}

func (a *QuarantineDeleteAction) Name() string { return "quarantine.delete" }

func (a *QuarantineDeleteAction) Execute(ctx context.Context, client *pmg.Client, input Data, params Params, filters map[string]string) (Data, error) {
	ids := getIDsFromInput(input)
	if len(ids) == 0 {
		return nil, fmt.Errorf("no mail IDs provided for deletion")
	}

	deleted := 0
	failed := 0

	for _, id := range ids {
		if err := client.DeleteQuarantine(id); err != nil {
			fmt.Printf("  [FAIL] Could not delete %s: %v\n", id, err)
			failed++
		} else {
			fmt.Printf("  [OK] Deleted: %s\n", id)
			deleted++
		}
	}

	return Data{
		"deleted": deleted,
		"failed":  failed,
		"total":   len(ids),
	}, nil
}

func getIDsFromInput(input Data) []string {
	if ids, ok := input["ids"]; ok {
		switch v := ids.(type) {
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

func init() {
	Register(&QuarantineListAction{})
	Register(&QuarantineDeliverAction{})
	Register(&QuarantineDeleteAction{})
}