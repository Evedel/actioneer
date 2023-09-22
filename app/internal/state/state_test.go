package state

import (
	"actioneer/internal/config"
	"testing"
)

func TestInitTemplateKeys(t *testing.T) {
	template := ""
	substitutionPrefix := "~"
	templateKeys := InitTemplateKeys(template, substitutionPrefix)
	if len(templateKeys) != 0 {
		t.Errorf("expected templateKeys to be empty, got %+v", templateKeys)
	}

	template = "foo"
	templateKeys = InitTemplateKeys(template, substitutionPrefix)
	if len(templateKeys) != 0 {
		t.Errorf("expected templateKeys to be empty, got %+v", templateKeys)
	}

	template = "~foo"
	templateKeys = InitTemplateKeys(template, substitutionPrefix)
	if len(templateKeys) != 1 {
		t.Errorf("expected templateKeys to be 1, got %+v", templateKeys)
	}
	if templateKeys[0] != "foo" {
		t.Errorf("expected templateKeys[0] to be foo, got %+v", templateKeys[0])
	}

	template = "~foo ~bar"
	templateKeys = InitTemplateKeys(template, substitutionPrefix)
	if len(templateKeys) != 2 {
		t.Errorf("expected templateKeys to be 2, got %+v", templateKeys)
	}
	if templateKeys[0] != "foo" {
		t.Errorf("expected templateKeys[0] to be foo, got %+v", templateKeys[0])
	}
	if templateKeys[1] != "bar" {
		t.Errorf("expected templateKeys[1] to be bar, got %+v", templateKeys[1])
	}

	template = "~foo bar baz"
	templateKeys = InitTemplateKeys(template, substitutionPrefix)
	if len(templateKeys) != 1 {
		t.Errorf("expected templateKeys to be 1, got %+v", templateKeys)
	}
	if templateKeys[0] != "foo" {
		t.Errorf("expected templateKeys[0] to be foo, got %+v", templateKeys[0])
	}
}

func TestInitState(t *testing.T) {
	config := config.Config{
		Version: "v1",
		Actions: []config.Action{
			{
				Alertname: "foo",
				Command:   "~foo",
			},
		},
		SubstitutionPrefix: "~",
	}

	state := InitState(config)
	if len(state.Actions) != 1 {
		t.Errorf("expected state.Actions to be 1, got %+v", state.Actions)
	}
	if len(state.Actions[0].TemplateKeys) != 1 {
		t.Errorf("expected state.Actions[0].TemplateKeys to be 1, got %+v", state.Actions[0].TemplateKeys)
	}
	if state.Actions[0].TemplateKeys[0] != "foo" {
		t.Errorf("expected state.Actions[0].TemplateKeys[0] to be foo, got %+v", state.Actions[0].TemplateKeys[0])
	}
	if state.Actions[0].Alertname != "foo" {
		t.Errorf("expected state.Actions[0].Alertname to be foo, got %+v", state.Actions[0].Alertname)
	}
	if state.Actions[0].CommandTemplate != "~foo" {
		t.Errorf("expected state.Actions[0].Command to be ~foo, got %+v", state.Actions[0].CommandTemplate)
	}
	if state.SubstitutionPrefix != "~" {
		t.Errorf("expected state.SubstitutionPrefix to be ~, got %+v", state.SubstitutionPrefix)
	}
}
