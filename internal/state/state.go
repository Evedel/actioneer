package state

import (
	"actioneer/internal/config"
	"strings"
)

type Action struct {
	Name 		  	string
	Alertname       string
	CommandTemplate string
	TemplateKeys    []string
}

type State struct {
	SubstitutionPrefix string
	AlertNameKey	   string
	Actions            []Action
}

func InitTemplateKeys(template string, substitutionPrefix string) (templateKeys []string) {
	for _, commandToken := range strings.Split(template, " ") {
		if strings.HasPrefix(commandToken, substitutionPrefix) {
			templateKeys = append(templateKeys, strings.TrimPrefix(commandToken, substitutionPrefix))
		}
	}
	return
}

func InitState(config config.Config) (state State) {
	state.SubstitutionPrefix = config.SubstitutionPrefix
	for _, action := range config.Actions {
		state.Actions = append(state.Actions, Action{
			Name: action.Name,
			Alertname:       action.Alertname,
			CommandTemplate: action.Command,
			TemplateKeys:    InitTemplateKeys(action.Command, config.SubstitutionPrefix),
		})
	}
	return
}

func (s State) GetActionByAlertName(alertname string) (action Action, found bool) {
	for _, action := range s.Actions {
		if action.Alertname == alertname {
			return action, true
		}
	}
	return action, false
}
