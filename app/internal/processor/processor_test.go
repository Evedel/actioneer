package processor

import (
	"actioneer/internal/logging"
	"actioneer/internal/state"
	th "actioneer/internal/testing_helper"
	"bytes"
	"testing"
)

func Test_ReadIncommingNotification_Ok(t *testing.T) {
	// given
	incommingBytes := []byte(`{"status":"firing","alerts":[{"status":"firing","labels":{"alertname":"High Pod Memory","pod":"test_pod_name","namespace":"monitoring"}}]}`)
	// when
	notification, err := ReadIncommingNotification(incommingBytes)
	// then
	if err != nil {
		t.Error("expected no error, got: " + err.Error())
	}
	th.AssertNil(t, err)
	th.AssertEqual(t, 1, len(notification.Alerts))
	th.AssertEqual(t, "firing", notification.Alerts[0].Status)
	th.AssertEqual(t, "High Pod Memory", notification.Alerts[0].Labels["alertname"])
	th.AssertEqual(t, "test_pod_name", notification.Alerts[0].Labels["pod"])
	th.AssertEqual(t, "monitoring", notification.Alerts[0].Labels["namespace"])
}

func Test_ReadIncommingNotification_Error(t *testing.T) {
	var buf bytes.Buffer
	logging.Init("error", &buf)

	// given
	incommingBytes := []byte(`{"status":"firing","alerts":[{"status":"firing","labels":{"alertname":"High Pod Memory","pod":"test_pod_name","namespace":"monitoring"}}`)
	// when
	_, err := ReadIncommingNotification(incommingBytes)
	// then
	th.AssertNotNil(t, err)
	th.AssertEqual(t, "unexpected end of JSON input", err.Error())
	th.AssertStringContains(t, "cannot unmarshal incomming bytes: {\\\"status\\\":\\\"firing\\\",\\\"alerts\\\":[{\\\"status\\\":\\\"firing\\\",\\\"labels\\\":{\\\"alertname\\\":\\\"High Pod Memory\\\",\\\"pod\\\":\\\"test_pod_name\\\",\\\"namespace\\\":\\\"monitoring\\\"}", buf.String())
}

func genAction(name string, alertName string, templateKeys []string) state.Action {
	if name == "" {
		name = "action1"
	}
	if alertName == "" {
		alertName = "High Pod Memory"
	}
	if len(templateKeys) == 0 {
		templateKeys = []string{"pod"}
	}
	return state.Action{
		Name:         name,
		Alertname:    alertName,
		TemplateKeys: templateKeys,
	}
}

func genState(actions []state.Action) state.State {
	if len(actions) == 0 {
		actions = []state.Action{
			genAction("", "", nil),
		}
	}
	return state.State{
		Actions: actions,
	}
}

func genAlert(status string, alertName string, pod string, namespace string) Alert {
	if status == "" {
		status = "firing"
	}
	if alertName == "" {
		alertName = "High Pod Memory"
	}
	if pod == "" {
		pod = "test_pod_name"
	}
	if namespace == "" {
		namespace = "monitoring"
	}
	return Alert{
		Status: status,
		Labels: map[string]string{
			"alertname": alertName,
			"pod":      pod,
			"namespace": namespace,
		},
	}
}

func Test_CheckActionNeeded_Ok(t *testing.T) {
	// given
	state := genState(nil)
	alert := genAlert("", "", "", "")
	// when
	actionNeeded := CheckActionNeeded(state, alert)
	// then
	th.AssertEqual(t, true, actionNeeded)
}

func Test_CheckActionNeeded_FoundButNotFiring(t *testing.T) {
	var buf bytes.Buffer
	logging.Init("debug", &buf)

	// given
	state := genState(nil)
	alert := genAlert("pending", "", "", "")

	// when
	actionNeeded := CheckActionNeeded(state, alert)
	// then
	th.AssertEqual(t, false, actionNeeded)
	th.AssertStringContains(t, "actions not found for alert=[High Pod Memory]", buf.String())
}

func Test_CheckActionNeeded_NotFound(t *testing.T) {
	var buf bytes.Buffer
	logging.Init("debug", &buf)

	// given
	state := genState(nil)
	alert := genAlert("", "High Pod CPU", "", "")

	// when
	actionNeeded := CheckActionNeeded(state, alert)
	// then
	th.AssertEqual(t, false, actionNeeded)
	th.AssertStringContains(t, "actions not found for alert=[High Pod CPU]", buf.String())
}

func Test_CheckTemplateLabelsPresent_Ok(t *testing.T) {
	// given
	action := genAction("", "", []string{"pod"})
	realLabelValues := map[string]string{
		"pod": "test_pod_name",
	}
	// when
	err := CheckTemplateLabelsPresent(action, realLabelValues)
	// then
	th.AssertNil(t, err)
}

func Test_CheckTemplateLabelsPresent_Error(t *testing.T) {
	var buf bytes.Buffer
	logging.Init("error", &buf)

	// given
	action := genAction("", "", []string{"pod"})
	realLabelValues := map[string]string{
		"namespace": "monitoring",
	}
	// when
	err := CheckTemplateLabelsPresent(action, realLabelValues)
	// then
	th.AssertNotNil(t, err)
	th.AssertEqual(t, "no label 'pod' were present on the alert, action=[action1] cannot be taken for alert=[High Pod Memory]", err.Error())
	th.AssertStringContains(t, "no label 'pod' were present on the alert, action=[action1] cannot be taken for alert=[High Pod Memory]", buf.String())
}

func Test_ExtractRealLabelValues_Ok(t *testing.T) {
	// given
	alert := genAlert("", "", "", "")
	// when
	realLabelValues := ExtractRealLabelValues(alert)
	// then
	th.AssertEqual(t, 3, len(realLabelValues))
	th.AssertEqual(t, "High Pod Memory", realLabelValues["alertname"])
	th.AssertEqual(t, "test_pod_name", realLabelValues["pod"])
	th.AssertEqual(t, "monitoring", realLabelValues["namespace"])
}
