package actions

import (
	"testing"
)

func TestGetIDsFromInput(t *testing.T) {
	tests := []struct {
		name  string
		input Data
		want  int
	}{
		{
			name:  "nil input",
			input: Data{},
			want:  0,
		},
		{
			name: "string slice",
			input: Data{
				"ids": []string{"1", "2", "3"},
			},
			want: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ids := getIDsFromInput(tt.input)
			if len(ids) != tt.want {
				t.Errorf("getIDsFromInput() = %d ids, want %d", len(ids), tt.want)
			}
		})
	}
}

func TestGetStringParam(t *testing.T) {
	tests := []struct {
		name   string
		params Params
		key    string
		want   string
	}{
		{"missing key", Params{}, "test", ""},
		{"string value", Params{"test": "hello"}, "test", "hello"},
		{"int value", Params{"test": 42}, "test", "42"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getStringParam(tt.params, tt.key)
			if got != tt.want {
				t.Errorf("getStringParam() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGetIntParam(t *testing.T) {
	tests := []struct {
		name   string
		params Params
		key    string
		want   int
	}{
		{"missing key", Params{}, "test", 0},
		{"int value", Params{"test": 42}, "test", 42},
		{"string value", Params{"test": "7"}, "test", 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getIntParam(tt.params, tt.key)
			if got != tt.want {
				t.Errorf("getIntParam() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestActionRegistry(t *testing.T) {
	allActions := All()
	required := []string{
		"quarantine.list",
		"quarantine.deliver",
		"quarantine.delete",
		"ruledb.who.list",
		"ruledb.who.add",
		"ruledb.who.remove",
		"transform.deduplicate",
		"transform.filter",
		"transform.extract",
		"report.console",
		"report.file",
		"report.json",
	}

	for _, name := range required {
		action, ok := allActions[name]
		if !ok {
			t.Errorf("action %q not registered", name)
		} else if action.Name() != name {
			t.Errorf("action Name() = %q, want %q", action.Name(), name)
		}
	}
}