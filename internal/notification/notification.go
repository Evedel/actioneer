package notification

import (
	"actioneer/internal/state"
	"fmt"
	"log/slog"
)

type Alert struct {
	Status string
	Name string
	Labels map[string]string
}

type Notification struct {
	Alerts []Alert
}

func ToInternal(ne NotificationExternal, state state.State) Notification {
	slog.Debug("Converting external notification to internal")

	n := Notification{}
	for _, ae := range ne.Alerts {
		alertName, ok := getAlertName(ae, state.AlertNameKey)
		if !ok {
			slog.Warn("no alert name label=[" + state.AlertNameKey + "], skipping=[" + fmt.Sprintf("%+v", ae)+"]")
			continue
		}
		n.Alerts = append(n.Alerts, Alert{
			Status: ae.Status,
			Name: alertName,
			Labels: ae.Labels,
		})
		n.Alerts[len(n.Alerts)-1].Labels["status"] = ae.Status
	}
	return n
}

func getAlertName(alert AlertExternal, alertNameKey string) (string, bool) {
	slog.Debug("Getting alert name from external alert")

	alertName := ""
	ok := false
	if alertName, ok = alert.Labels[alertNameKey]; !ok {
		slog.Debug("no alert name label=[" + alertNameKey + "], skipping=[" + fmt.Sprintf("%+v", alert)+"]")
	}
	return alertName, ok
}
