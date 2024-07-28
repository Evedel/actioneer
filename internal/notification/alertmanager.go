package notification

import (
	"encoding/json"
	"fmt"
	"log/slog"
)

type AlertExternal struct {
	Status string
	Labels map[string]string
}

type NotificationExternal struct {
	Alerts []AlertExternal
}

func ReadAlertmanagerNotification(bytes []byte) (NotificationExternal, error) {
	slog.Debug("incomming bytes: " + fmt.Sprintf("%+v", string(bytes)))
	var ne NotificationExternal
	err := json.Unmarshal(bytes, &ne)
	if err != nil {
		slog.Error("cannot unmarshal incomming bytes: " + fmt.Sprintf("%+v", string(bytes)))
		slog.Error(err.Error())
	}
	return ne, err
}
