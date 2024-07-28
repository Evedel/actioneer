package notification

import (
	"actioneer/internal/logging"
	th "actioneer/internal/testing_helper"
	"bytes"
	"testing"
)

func gen_AlertExternal(status string, labels th.Dict) AlertExternal {
	if status == "" {
		status = "firing"
	}
	if labels == nil {
		labels = th.Dict{
			"alertname": "High Pod Memory",
			"pod":      "test_pod_name",
			"namespace": "monitoring",
		}
	}
	return AlertExternal{
		Status: status,
		Labels: labels,
	}
}

func get_NotificationExternal(statuses []string, labels []th.Dict) (NotificationExternal) {
	ne := NotificationExternal{}
	for i, label := range labels {
		ae := gen_AlertExternal(statuses[i], label)
		ne.Alerts = append(ne.Alerts, ae)
	}
	return ne
}

func Test_getAlertName_Ok(t *testing.T) {
	alert := gen_AlertExternal("", th.Dict{"alertname": "High Pod Memory"})
	alertNameKey := "alertname"
	expected := "High Pod Memory"

	actual, ok := getAlertName(alert, alertNameKey)
	th.AssertEqual(t, expected, actual)
	th.AssertEqual(t, true, ok)
}

func Test_getAlertName_NoAlertNameKey(t *testing.T) {
	var buf bytes.Buffer
	logging.Init("debug", &buf)

	alert := gen_AlertExternal("", th.Dict{"myalertname": "High Pod Memory"})

	alertNameKey := "alertname"
	expected := ""

	actual, ok := getAlertName(alert, alertNameKey)
	th.AssertEqual(t, expected, actual)
	th.AssertEqual(t, false, ok)
	th.AssertStringContains(t, "no alert name label=[alertname], skipping=[{Status:firing Labels:map[myalertname:High Pod Memory", buf.String())
}

func Test_ToInternal_NameAdded(t *testing.T) {
	statuses := []string{"firing", "resolved"}
	labels := []th.Dict{
		{"alertname": "High Pod Memory"},
		{"alertname": "Low Pod Memory"},
	}
	ne := get_NotificationExternal(statuses, labels)

	state := th.GetState(th.Dict{"alertname": "alertname"})

	n_actual := ToInternal(ne, state)
	th.AssertEqual(t, 2, len(n_actual.Alerts))
	th.AssertEqual(t, "High Pod Memory", n_actual.Alerts[0].Name)
	th.AssertEqual(t, "Low Pod Memory", n_actual.Alerts[1].Name)
}

func Test_ToInternal_NoAlertNameKey_Skipped(t *testing.T) {
	var buf bytes.Buffer
	logging.Init("debug", &buf)

	statuses := []string{"firing", "resolved"}
	labels := []th.Dict{
		{"alertname": "High Pod Memory"},
		{"myalertname": "Low Pod Memory"},
	}
	ne := get_NotificationExternal(statuses, labels)

	state := th.GetState(th.Dict{"alertname": "alertname"})

	n_actual := ToInternal(ne, state)
	th.AssertEqual(t, 1, len(n_actual.Alerts))
	th.AssertEqual(t, "High Pod Memory", n_actual.Alerts[0].Name)
	th.AssertStringContains(t, "no alert name label=[alertname], skipping=[{Status:resolved Labels:map[myalertname:Low Pod Memory", buf.String())
}

func Test_ToInternal_StatusAdded(t *testing.T) {
	statuses := []string{"test_status"}
	labels := []th.Dict{
		{"alertname": "High Pod Memory"},
	}
	ne := get_NotificationExternal(statuses, labels)

	state := th.GetState(th.Dict{"alertname": "alertname"})

	n_actual := ToInternal(ne, state)
	th.AssertEqual(t, 1, len(n_actual.Alerts))
	th.AssertEqual(t, "test_status", n_actual.Alerts[0].Status)
}
