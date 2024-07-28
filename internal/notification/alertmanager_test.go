package notification

import (
	"actioneer/internal/logging"
	th "actioneer/internal/testing_helper"
	"bytes"
	"testing"
)

func Test_ReadAlertmanagerNotification_Ok(t *testing.T) {
	// given
	incommingBytes := []byte(`{"status":"firing","alerts":[{"status":"firing","labels":{"alertname":"High Pod Memory","pod":"test_pod_name","namespace":"monitoring"}}]}`)
	// when
	notification, err := ReadAlertmanagerNotification(incommingBytes)
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

func Test_ReadAlertmanagerNotification_CustomAlertNameKey_Ok(t *testing.T) {
	// given
	incommingBytes := []byte(`{"status":"firing","alerts":[{"status":"firing","labels":{"myalertname":"High Pod Memory","pod":"test_pod_name","namespace":"monitoring"}}]}`)
	// when
	notification, err := ReadAlertmanagerNotification(incommingBytes)
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

func Test_ReadAlertmanagerNotification_Error(t *testing.T) {
	var buf bytes.Buffer
	logging.Init("error", &buf)

	// given
	incommingBytes := []byte(`{"status":"firing","alerts":[{"status":"firing","labels":{"alertname":"High Pod Memory","pod":"test_pod_name","namespace":"monitoring"}}`)
	// when
	_, err := ReadAlertmanagerNotification(incommingBytes)
	// then
	th.AssertNotNil(t, err)
	th.AssertEqual(t, "unexpected end of JSON input", err.Error())
	th.AssertStringContains(t, "cannot unmarshal incomming bytes: {\\\"status\\\":\\\"firing\\\",\\\"alerts\\\":[{\\\"status\\\":\\\"firing\\\",\\\"labels\\\":{\\\"alertname\\\":\\\"High Pod Memory\\\",\\\"pod\\\":\\\"test_pod_name\\\",\\\"namespace\\\":\\\"monitoring\\\"}", buf.String())
}
