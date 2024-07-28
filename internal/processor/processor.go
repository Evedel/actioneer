package processor

import (
	"actioneer/internal/command"
	"actioneer/internal/notification"
	"actioneer/internal/state"
	"fmt"
	"log/slog"
	"strings"
)

func CheckActionNeeded(state state.State, alert notification.Alert) bool {
	_, found := state.GetActionByAlertName(alert.Name)
	if found && (alert.Status == "firing") {
		slog.Debug("action found for alert=[" + fmt.Sprint(alert.Name)+"]")
		return true
	}
	slog.Debug("actions not found for alert=[" + fmt.Sprint(alert.Name)+"]")
	return false
}

func CheckTemplateLabelsPresent(action state.Action, realLabelValues map[string]string) error {
	for _, templateKey := range action.TemplateKeys {
		if _, ok := realLabelValues[templateKey]; !ok {
			errString := "no label '" + templateKey + "' were present on the alert, action=[" + fmt.Sprintf("%+v", action.Name) + "] cannot be taken for alert=[" + fmt.Sprintf("%+v", action.Alertname)+"]"
			slog.Error(errString)
			err := fmt.Errorf(errString)
			return err
		}
	}
	return nil
}

func ExtractRealLabelValues(alert notification.Alert) (map[string]string) {
	realLabelValues := make(map[string]string)
	for k, v := range alert.Labels {
		realLabelValues[k] = v
	}
	return realLabelValues
}

func CompileCommandTemplate(action state.Action, realLabelValues map[string]string, substitutionPrefix string) string {
	commandReady := action.CommandTemplate
	for k, v := range realLabelValues {
		commandReady = strings.ReplaceAll(commandReady, substitutionPrefix+k, v)
	}
	return commandReady
}

func TakeActions(shell command.ICommandRunner, state state.State, notification notification.Notification, isDryRun bool) error {
	slog.Debug("incomming notification=[" + fmt.Sprint(notification)+"]")

	if len(notification.Alerts) == 0 {
		slog.Error("no alerts in notification=[" + fmt.Sprint(notification)+"]")
		return nil
	}

	for _, alert := range notification.Alerts {
		if !CheckActionNeeded(state, alert) {
			continue
		}

		action, _ := state.GetActionByAlertName(alert.Name)
		slog.Debug("command template=[" + fmt.Sprint(action.CommandTemplate)+"]")

		realLabelValues := ExtractRealLabelValues(alert)
		slog.Debug("found lables on the real alert=[" + fmt.Sprint(realLabelValues)+"]")

		err := CheckTemplateLabelsPresent(action, realLabelValues)
		if err != nil {
			slog.Error(err.Error())
			return err
		}

		commandReady := CompileCommandTemplate(action, realLabelValues, state.SubstitutionPrefix)

		command.Execute(shell, commandReady, isDryRun)
	}
	return nil
}
