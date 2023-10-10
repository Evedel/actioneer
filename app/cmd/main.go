package main

import (
	"actioneer/internal/args"
	"actioneer/internal/command"
	"actioneer/internal/config"
	"actioneer/internal/logging"
	"actioneer/internal/state"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

var alertNameLabel = "alertname"

type Alert struct {
	Status string
	Labels map[string]string
}

type Notification struct {
	Alerts []Alert
}

type Server struct {
	IsDryRun bool
	State    state.State
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	slog.Debug("incomming request body: " + fmt.Sprintf("%+v", string(bytes)))
	var notification Notification
	if err := json.Unmarshal(bytes, &notification); err != nil {
		slog.Error("cannot unmarshal incomming notification: " + fmt.Sprintf("%+v", string(bytes)))
		slog.Error(err.Error())
	}

	slog.Debug("processing notification: " + fmt.Sprint(notification))
	if len(notification.Alerts) == 0 {
		slog.Error("no alerts in notification: " + fmt.Sprint(notification))
		return
	}

	for _, alert := range notification.Alerts {
		if alertName, ok := alert.Labels[alertNameLabel]; ok {
			slog.Debug("processing alert: " + alertName)

			action, err := s.State.GetActionByAlertName(alertName)
			if err != nil {
				slog.Warn(err.Error())
				return
			}
			slog.Debug("command template: " + fmt.Sprint(action.CommandTemplate))

			labelValues := make(map[string]string)
			for _, templateKey := range action.TemplateKeys {
				if value, ok := alert.Labels[templateKey]; ok {
					labelValues[templateKey] = value
				} else {
					slog.Error("no label '" + templateKey + "' in alert, skipping alert: " + fmt.Sprintf("%+v", alert) + " and action: " + fmt.Sprintf("%+v", action))
					return
				}
			}

			commandReady := action.CommandTemplate
			for k, v := range labelValues {
				commandReady = strings.ReplaceAll(commandReady, s.State.SubstitutionPrefix+k, v)
			}

			command.Execute(command.CommandRunner{}, commandReady, s.IsDryRun)
		} else {
			slog.Error("no alert name label '" + alertNameLabel + "', skipping: " + fmt.Sprintf("%+v", alert))
		}
	}
}

func main() {
	args := args.Parse()

	if err := logging.Init(*args.LogLevel, os.Stdout); err != nil {
		os.Exit(2)
	}

	cfg, err := config.Read(config.ConfigReader{}, *args.ConfigPath)
	if err != nil {
		os.Exit(2)
	}
	if !config.IsValid(cfg) {
		os.Exit(2)
	}

	actions := state.InitState(cfg)

	s := Server{IsDryRun: *args.IsDryRun, State: actions}
	mux := http.NewServeMux()
	mux.Handle("/", s)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		slog.Error(err.Error())
		os.Exit(2)
	}
}
