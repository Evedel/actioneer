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

func TestInitState_NonDefaultSubstitutionPrefix_NotFound(t *testing.T) {
	config1 := config.Config{
		Version: "v1",
		Actions: []config.Action{
			{
				Alertname: "foo",
				Command:   "~foo",
			},
		},
		SubstitutionPrefix: "$label.",
	}

	state1 := InitState(config1)
	if len(state1.Actions) != 1 {
		t.Errorf("expected state.Actions to be 1, got %+v", state1.Actions)
	}
	if len(state1.Actions[0].TemplateKeys) != 0 {
		t.Errorf("expected state.Actions[0].TemplateKeys to be empty, got %+v", state1.Actions[0].TemplateKeys)
	}
	if state1.SubstitutionPrefix != "$label." {
		t.Errorf("expected state.SubstitutionPrefix to be \"$label.\", got %+v", state1.SubstitutionPrefix)
	}
}

func TestInitState_NonDefaultSubstitutionPrefix_Found(t *testing.T) {
	config1 := config.Config{
		Version: "v1",
		Actions: []config.Action{
			{
				Alertname: "foo",
				Command:   "~foo $label.bar",
			},
		},
		SubstitutionPrefix: "$label.",
	}

	state2 := InitState(config1)
	if len(state2.Actions) != 1 {
		t.Errorf("expected state.Actions to be 1, got %+v", state2.Actions)
	}
	if len(state2.Actions[0].TemplateKeys) != 1 {
		t.Errorf("expected state.Actions[0].TemplateKeys to be 1, got %+v", state2.Actions[0].TemplateKeys)
	}
	if state2.Actions[0].TemplateKeys[0] != "bar" {
		t.Errorf("expected state.Actions[0].TemplateKeys[0] to be \"bar\", got %+v", state2.Actions[0].TemplateKeys[0])
	}
	if state2.Actions[0].Alertname != "foo" {
		t.Errorf("expected state.Actions[0].Alertname to be foo, got %+v", state2.Actions[0].Alertname)
	}
	if state2.Actions[0].CommandTemplate != "~foo $label.bar" {
		t.Errorf("expected state.Actions[0].Command to be \"~foo $label.bar\", got %+v", state2.Actions[0].CommandTemplate)
	}
	if state2.SubstitutionPrefix != "$label." {
		t.Errorf("expected state.SubstitutionPrefix to be \"$label.bar\", got %+v", state2.SubstitutionPrefix)
	}
}

func TestGetActionByAlertName_NotFond(t *testing.T) {
	config1 := config.Config{
		Version: "v1",
		Actions: []config.Action{
			{
				Alertname: "foo",
				Command:   "~foo $label.bar",
			}, {
				Alertname: "bar",
				Command:   "~bar $label.foo",
			}, {
				Alertname: "baz",
				Command:   "~baz $label.baz",
			},
		},
		SubstitutionPrefix: "$label.",
	}

	state1 := InitState(config1)
	action, err := state1.GetActionByAlertName("qux")
	if err == nil {
		t.Errorf("expected err to be not nil, got %+v", err)
	}
	if action.Alertname != "" {
		t.Errorf("expected action.Alertname to be empty, got %+v", action.Alertname)
	}
	if action.CommandTemplate != "" {
		t.Errorf("expected action.CommandTemplate to be empty, got %+v", action.CommandTemplate)
	}
	if len(action.TemplateKeys) != 0 {
		t.Errorf("expected action.TemplateKeys to be empty, got %+v", action.TemplateKeys)
	}
}

func TestGetActionByAlertName_Fond(t *testing.T) {
	config1 := config.Config{
		Version: "v1",
		Actions: []config.Action{
			{
				Alertname: "foo",
				Command:   "~foo $label.bar",
			}, {
				Alertname: "bar",
				Command:   "~bar $label.foo",
			}, {
				Alertname: "baz",
				Command:   "~baz $label.baz",
			},
		},
		SubstitutionPrefix: "$label.",
	}

	state1 := InitState(config1)
	action, err := state1.GetActionByAlertName("bar")
	if err != nil {
		t.Errorf("expected err to be nil, got %+v", err)
	}
	if action.Alertname != "bar" {
		t.Errorf("expected action.Alertname to be bar, got %+v", action.Alertname)
	}
	if action.CommandTemplate != "~bar $label.foo" {
		t.Errorf("expected action.CommandTemplate to be \"~bar $label.foo\", got %+v", action.CommandTemplate)
	}
	if len(action.TemplateKeys) != 1 {
		t.Errorf("expected action.TemplateKeys to be 1, got %+v", action.TemplateKeys)
	}
	if action.TemplateKeys[0] != "foo" {
		t.Errorf("expected action.TemplateKeys[0] to be foo, got %+v", action.TemplateKeys[0])
	}
}
