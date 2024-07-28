package processor

import (
	"actioneer/internal/command"
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

func Test_ReadIncommingNotification_CustomAlertNameKey_Ok(t *testing.T) {
	// given
	incommingBytes := []byte(`{"status":"firing","alerts":[{"status":"firing","labels":{"myalertname":"High Pod Memory","pod":"test_pod_name","namespace":"monitoring"}}]}`)
	// when
	notification, err := ReadIncommingNotification(incommingBytes)
	// then
	if err != nil {
		t.Error("expected no error, got: " + err.Error())
	}
	th.AssertNil(t, err)
	th.AssertEqual(t, 1, len(notification.Alerts))
	th.AssertEqual(t, "firing", notification.Alerts[0].Status)
	th.AssertEqual(t, "High Pod Memory", notification.Alerts[0].Labels["myalertname"])
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

func genAction(name string, alertName string, commandTemplate string, templateKeys []string) state.Action {
	if name == "" {
		name = "action1"
	}
	if alertName == "" {
		alertName = "High Pod Memory"
	}
	if commandTemplate == "" {
		commandTemplate = "echo $pod"
	}
	if len(templateKeys) == 0 {
		templateKeys = []string{"pod"}
	}
	return state.Action{
		Name:         		name,
		Alertname:    		alertName,
		CommandTemplate: 	commandTemplate,
		TemplateKeys: 		templateKeys,
	}
}

func genState(actions []state.Action) state.State {
	if len(actions) == 0 {
		actions = []state.Action{
			genAction("", "", "", nil),
		}
	}
	return state.State{
		Actions: actions,
		SubstitutionPrefix: "$",
	}
}

func genAlert(status string, labels map[string]string) Alert {
	if status == "" {
		status = "firing"
	}
	if labels == nil {
		labels = map[string]string{
			"alertname": "High Pod Memory",
			"pod":      "test_pod_name",
			"namespace": "monitoring",
		}
	}
	return Alert{
		Status: status,
		Labels: labels,
	}
}

func genNotification(alerts []Alert) Notification {
	if alerts == nil {
		alerts = []Alert{
			genAlert("", nil),
		}
	}
	return Notification{
		Alerts: alerts,
	}
}

func Test_CheckActionNeeded_Ok(t *testing.T) {
	// given
	state := genState(nil)
	alert := genAlert("", nil)
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
	alert := genAlert("pending", nil)

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
	alert := genAlert("", map[string]string{"alertname" : "High Pod CPU"})

	// when
	actionNeeded := CheckActionNeeded(state, alert)
	// then
	th.AssertEqual(t, false, actionNeeded)
	th.AssertStringContains(t, "actions not found for alert=[High Pod CPU]", buf.String())
}

func Test_CheckTemplateLabelsPresent_Ok(t *testing.T) {
	// given
	action := genAction("", "", "", []string{"pod"})
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
	action := genAction("", "", "", []string{"pod"})
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
	alert := genAlert("", nil)
	// when
	realLabelValues := ExtractRealLabelValues(alert)
	// then
	th.AssertEqual(t, 3, len(realLabelValues))
	th.AssertEqual(t, "High Pod Memory", realLabelValues["alertname"])
	th.AssertEqual(t, "test_pod_name", realLabelValues["pod"])
	th.AssertEqual(t, "monitoring", realLabelValues["namespace"])
}

func Test_CompileCommandTemplate_Ok(t *testing.T) {
	// given
	action := genAction("", "", "test delete TFDpod1 TFDnamespace1", nil)
	realLabelValues := map[string]string{
		"pod1": "test_pod_name",
		"namespace1": "monitoring",
	}
	// when
	commandReady := CompileCommandTemplate(action, realLabelValues, "TFD")
	// then
	th.AssertEqual(t, "test delete test_pod_name monitoring", commandReady)
}

func Test_TakeActions_Ok(t *testing.T) {
	// given
	shell := command.FakeCommandRunner{}
	state := genState(
		[]state.Action{
			genAction("action1", "High Pod Memory", "echo $pod", []string{"pod"}),
			genAction("action2", "High Pod CPU", "echo $pod $namespace", []string{"pod", "namespace"}),
			genAction("action3", "High Pod Storage", "echo $cluster", []string{"pod", "cluster"}),
		},
	)
	notification := genNotification(
		[]Alert{
			genAlert("firing", map[string]string{"alertname": "High Pod Memory", "pod": "test_pod_name", "namespace":"monitoring"}),
			genAlert("pending", map[string]string{"alertname": "High Pod CPU", "pod": "test_pod_name", "namespace":"monitoring"}),
			genAlert("resolved", map[string]string{"alertname": "High Pod Storage", "pod": "test_pod_name", "cluster":"test"}),
		},
	)
	// when
	err := TakeActions(&shell, state, notification, false)

	// then
	th.AssertNil(t, err)
	th.AssertEqual(t, 1, len(shell.Calls))
	th.AssertEqual(t, "bash -c echo test_pod_name", shell.Calls[0])
}

func Test_TakeActions_NoAlerts(t *testing.T) {
	var buf bytes.Buffer
	logging.Init("error", &buf)

	// given
	shell := command.FakeCommandRunner{}
	state := genState(
		[]state.Action{
			genAction("action1", "High Pod Memory", "echo $pod", []string{"pod"}),
			genAction("action2", "High Pod CPU", "echo $pod $namespace", []string{"pod", "namespace"}),
			genAction("action3", "High Pod Storage", "echo $cluster", []string{"pod", "cluster"}),
		},
	)
	notification := genNotification(
		[]Alert{},
	)
	// when
	err := TakeActions(&shell, state, notification, false)
	// then
	th.AssertNil(t, err)
	th.AssertEqual(t, 0, len(shell.Calls))
	th.AssertStringContains(t, "no alerts in notification", buf.String())
}

func Test_TakeActions_NoAlertName(t *testing.T) {
	var buf bytes.Buffer
	logging.Init("error", &buf)

	// given
	shell := command.FakeCommandRunner{}
	state := genState(
		[]state.Action{
			genAction("action1", "High Pod Memory", "echo $pod", []string{"pod"}),
			genAction("action2", "High Pod CPU", "echo $pod $namespace", []string{"pod", "namespace"}),
			genAction("action3", "High Pod Storage", "echo $cluster", []string{"pod", "cluster"}),
		},
	)
	notification := genNotification(
		[]Alert{
			genAlert("firing", map[string]string{"pod": "test_pod_name", "namespace":"monitoring"}),
		},
	)
	// when
	err := TakeActions(&shell, state, notification, false)
	// then
	th.AssertNil(t, err)
	th.AssertEqual(t, 0, len(shell.Calls))
	th.AssertStringContains(t, "no alert name label=[alertname], skipping=[{Status:firing Labels:map[namespace:monitoring pod:test_pod_name]}]", buf.String())
}

func Test_TakeActions_NoActionFound(t *testing.T) {
	var buf bytes.Buffer
	logging.Init("debug", &buf)

	// given
	shell := command.FakeCommandRunner{}
	state := genState(
		[]state.Action{
			genAction("action1", "High Pod Memory", "echo $pod", []string{"pod"}),
			genAction("action2", "High Pod CPU", "echo $pod $namespace", []string{"pod", "namespace"}),
			genAction("action3", "High Pod Storage", "echo $cluster", []string{"pod", "cluster"}),
		},
	)
	notification := genNotification(
		[]Alert{
			genAlert("firing", map[string]string{"alertname": "High Pod TTL", "pod": "test_pod_name", "namespace":"monitoring", "cluster":"test"}),
		},
	)
	// when
	err := TakeActions(&shell, state, notification, false)
	// then
	th.AssertNil(t, err)
	th.AssertEqual(t, 0, len(shell.Calls))
	th.AssertStringContains(t, "actions not found for alert=[High Pod TTL]", buf.String())
}

func Test_TakeActions_NoLabel(t *testing.T) {
	var buf bytes.Buffer
	logging.Init("error", &buf)

	// given
	shell := command.FakeCommandRunner{}
	state := genState(
		[]state.Action{
			genAction("action1", "High Pod Memory", "echo $pod", []string{"pod"}),
			genAction("action2", "High Pod CPU", "echo $pod $namespace", []string{"pod", "namespace"}),
			genAction("action3", "High Pod Storage", "echo $cluster", []string{"pod", "cluster"}),
		},
	)
	notification := genNotification(
		[]Alert{
			genAlert("firing", map[string]string{"alertname": "High Pod Memory", "namespace":"monitoring", "cluster":"test"}),
		},
	)
	// when
	err := TakeActions(&shell, state, notification, false)
	// then
	th.AssertNotNil(t, err)
	th.AssertEqual(t, "no label 'pod' were present on the alert, action=[action1] cannot be taken for alert=[High Pod Memory]", err.Error())
	th.AssertEqual(t, 0, len(shell.Calls))
	th.AssertStringContains(t, "no label 'pod' were present on the alert, action=[action1] cannot be taken for alert=[High Pod Memory]", buf.String())
}
