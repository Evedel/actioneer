package main

import (
	"actioneer/internal/args"
	"actioneer/internal/config"
	"actioneer/internal/logging"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
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
	Config   config.Config
	IsDryRun bool
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
	if len(notification.Alerts) > 0 {
		for _, alert := range notification.Alerts {
			if alertName, ok := alert.Labels[alertNameLabel]; ok {
				slog.Debug("processing alert: " + alertName)
				for _, action := range s.Config.Actions {
					if action.Alertname == alertName {
						slog.Debug("command template: " + fmt.Sprint(action.Command))

						labelValues := make(map[string]string)
						for _, cmdWord := range strings.Split(action.Command, " ") {
							if strings.HasPrefix(cmdWord, "~") {
								if value, ok := alert.Labels[strings.TrimPrefix(cmdWord, "~")]; ok {
									labelValues[cmdWord] = value
								} else {
									slog.Error("no label '" + strings.TrimPrefix(cmdWord, "~") + "' in alert, skipping alert: " + fmt.Sprintf("%+v", alert) + " and action: " + fmt.Sprintf("%+v", action))
									return
								}
							}
						}

						for k, v := range labelValues {
							action.Command = strings.ReplaceAll(action.Command, k, v)
						}

						slog.Debug("processing command: " + fmt.Sprint(action.Command))
						if s.IsDryRun {
							slog.Info("dry run: " + action.Command)
						} else {
							cmd := exec.Command("bash", "-c", action.Command)
							var stdout strings.Builder
							var stderr strings.Builder
							cmd.Stdout = &stdout
							cmd.Stderr = &stderr
							if err := cmd.Run(); err != nil {
								slog.Error("command execution failed: " + action.Command)
								slog.Error(err.Error())
							}
							slog.Info("stdout: " + stdout.String())
							if stderr.String() == "" {
								slog.Info("stderr: " + stderr.String())
							} else {
								slog.Error("stderr: " + stderr.String())
							}
						}
					}
				}
			} else {
				slog.Error("no alert name label '" + alertNameLabel + "', skipping: " + fmt.Sprintf("%+v", alert))
			}
		}
	}
}

func main() {
	args := args.Parse()

	if err := logging.Init(*args.LogLevel, nil); err != nil {
		os.Exit(2)
	}

	cfg, err := config.Read(config.ConfigReader{}, *args.ConfigPath)
	if err != nil {
		os.Exit(2)
	}
	if !config.IsValid(cfg) {
		os.Exit(2)
	}

	mux := http.NewServeMux()
	s := Server{Config: cfg, IsDryRun: *args.IsDryRun}
	mux.Handle("/", s)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		slog.Error(err.Error())
		os.Exit(2)
	}
}
