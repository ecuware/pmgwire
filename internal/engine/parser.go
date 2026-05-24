package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/ecuware/pmgwire/internal/actions"
	"gopkg.in/yaml.v3"
)

type Workflow struct {
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	Version     string            `yaml:"version"`
	Auth        AuthConfig        `yaml:"auth"`
	Vars        map[string]VarDef `yaml:"vars"`
	Steps       []Step            `yaml:"steps"`
}

type AuthConfig struct {
	Host     string `yaml:"host"`
	Token    string `yaml:"token"`
	Insecure bool   `yaml:"insecure"`
}

type VarDef struct {
	Default  interface{} `yaml:"default"`
	Prompt   string      `yaml:"prompt"`
	Required bool        `yaml:"required"`
}

type Step struct {
	ID         string                 `yaml:"id"`
	Action     string                 `yaml:"action"`
	Params     map[string]interface{} `yaml:"params"`
	Filters    map[string]string      `yaml:"filters"`
	Input      string                 `yaml:"input"`
	Output     string                 `yaml:"output"`
	Confirm    bool                   `yaml:"confirm"`
	OnError    string                 `yaml:"on_error"`
	RetryCount int                    `yaml:"retry_count"`
}

func ParseWorkflow(path string) (*Workflow, error) {
	cleanPath := filepath.Clean(path)
	if !filepath.IsAbs(cleanPath) {
		absPath, err := filepath.Abs(cleanPath)
		if err != nil {
			return nil, fmt.Errorf("resolving workflow path: %w", err)
		}
		cleanPath = absPath
	}

	data, err := os.ReadFile(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("reading workflow file: %w", err)
	}

	var wf Workflow
	if err := yaml.Unmarshal(data, &wf); err != nil {
		return nil, fmt.Errorf("parsing workflow YAML: %w", err)
	}

	if wf.Name == "" {
		return nil, fmt.Errorf("workflow name is required")
	}
	if wf.Version == "" {
		wf.Version = "1.0"
	}
	if len(wf.Steps) == 0 {
		return nil, fmt.Errorf("workflow must have at least one step")
	}

	stepIDs := make(map[string]bool)
	for i := range wf.Steps {
		step := &wf.Steps[i]
		if step.ID == "" {
			return nil, fmt.Errorf("each step must have an id")
		}
		if step.Action == "" {
			return nil, fmt.Errorf("step %s must have an action", step.ID)
		}
		if stepIDs[step.ID] {
			return nil, fmt.Errorf("duplicate step id: %s", step.ID)
		}
		stepIDs[step.ID] = true

		if step.OnError == "" {
			step.OnError = "stop"
		}
		if step.RetryCount == 0 && step.OnError == "retry" {
			step.RetryCount = 3
		}
	}

	if wf.Auth.Host == "" {
		wf.Auth.Host = "https://localhost:8006"
	}

	return &wf, nil
}

func ResolveVars(wf *Workflow) (map[string]interface{}, error) {
	resolved := make(map[string]interface{})

	for name, def := range wf.Vars {
		envKey := "PMGWIRE_" + strings.ToUpper(strings.ReplaceAll(name, "-", "_"))
		if envVal := os.Getenv(envKey); envVal != "" {
			resolved[name] = envVal
			continue
		}

		if def.Prompt != "" {
			fmt.Printf("%s (default: %v): ", def.Prompt, def.Default)
			var input string
			fmt.Scanln(&input)
			if input != "" {
				resolved[name] = input
			} else if def.Default != nil {
				resolved[name] = def.Default
			} else if def.Required {
				return nil, fmt.Errorf("required variable %s is missing", name)
			}
		} else if def.Default != nil {
			resolved[name] = def.Default
		} else if def.Required {
			return nil, fmt.Errorf("required variable %s is missing", name)
		}
	}

	return resolved, nil
}

func ResolveTemplates(wf *Workflow, vars map[string]interface{}, stepOutputs map[string]map[string]interface{}) error {
	tmplData := map[string]interface{}{
		"vars":  vars,
		"steps": stepOutputs,
	}

	for name, val := range vars {
		tmplData[name] = val
	}

	for i := range os.Environ() {
		entry := os.Environ()[i]
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) == 2 && strings.HasPrefix(parts[0], "PMG_") {
			tmplData[parts[0]] = parts[1]
		}
	}

	for i := range wf.Steps {
		step := &wf.Steps[i]

		for k, v := range step.Params {
			resolved, err := resolveTemplate(v, tmplData)
			if err != nil {
				return fmt.Errorf("resolving template in step %s param %s: %w", step.ID, k, err)
			}
			step.Params[k] = resolved
		}

		for k, v := range step.Filters {
			resolved, err := resolveTemplate(v, tmplData)
			if err != nil {
				return fmt.Errorf("resolving template in step %s filter %s: %w", step.ID, k, err)
			}
			if str, ok := resolved.(string); ok {
				step.Filters[k] = str
			}
		}
	}

	return nil
}

func ResolveStepTemplates(step *Step, vars map[string]interface{}, stepOutputs map[string]actions.Data) error {
	tmplData := map[string]interface{}{
		"vars":  vars,
		"steps": stepOutputs,
	}

	for name, val := range vars {
		tmplData[name] = val
	}

	for k, v := range step.Params {
		resolved, err := resolveTemplate(v, tmplData)
		if err != nil {
			return fmt.Errorf("resolving template in step %s param %s: %w", step.ID, k, err)
		}
		step.Params[k] = resolved
	}

	for k, v := range step.Filters {
		resolved, err := resolveTemplate(v, tmplData)
		if err != nil {
			return fmt.Errorf("resolving template in step %s filter %s: %w", step.ID, k, err)
		}
		if str, ok := resolved.(string); ok {
			step.Filters[k] = str
		}
	}

	if strings.Contains(step.Input, "{{") {
		resolved, err := resolveTemplate(step.Input, tmplData)
		if err != nil {
			return fmt.Errorf("resolving template in step %s input: %w", step.ID, err)
		}
		if str, ok := resolved.(string); ok {
			step.Input = str
		}
	}

	return nil
}

func resolveTemplate(val interface{}, data map[string]interface{}) (interface{}, error) {
	str, ok := val.(string)
	if !ok {
		return val, nil
	}

	if !strings.Contains(str, "{{") {
		return val, nil
	}

	tmpl, err := template.New("").Option("missingkey=error").Parse(str)
	if err != nil {
		return nil, fmt.Errorf("parsing template: %w", err)
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("executing template: %w", err)
	}

	return buf.String(), nil
}