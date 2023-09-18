package main

import (
	"actioneer/internal/args"
	"actioneer/internal/logging"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v3"
)

var alertNameLabel = "alertname"

type Action struct {
	Alertname string
	Command   string
}

type Config struct {
	Version string
	Actions []Action
}

type Alert struct {
	Status string
	Labels map[string]string
}

type Notification struct {
	Alerts []Alert
}

type Server struct {
	Config   Config
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

func ReadConfig(path string) (config Config) {
	file, err := os.Open(path)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(2)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(2)
	}

	if err := yaml.Unmarshal(bytes, &config); err != nil {
		slog.Error(err.Error())
		os.Exit(2)
	}

	return config
}

func ValidateConfig(config Config) (result bool) {
	if config.Version != "v1" {
		slog.Error("wrong config version: " + config.Version)
		return false
	}

	if len(config.Actions) == 0 {
		slog.Error("no actions defined")
		return false
	}

	for _, action := range config.Actions {
		if action.Alertname == "" {
			slog.Error("empty alertname in action: " + fmt.Sprint(action))
			return false
		}
		if action.Command == "" {
			slog.Error("empty command in action: " + fmt.Sprint(action))
			return false
		}
	}

	for i, action := range config.Actions {
		for j, action2 := range config.Actions {
			if i != j && action.Alertname == action2.Alertname {
				slog.Error("duplicate alertname in actions: " + fmt.Sprint(action) + " and " + fmt.Sprint(action2))
				return false
			}
		}
	}

	slog.Debug("config: " + fmt.Sprintf("%+v", config))

	return true
}

func main() {
	args := args.Parse()

	if err := logging.Init(*args.LogLevel, nil); err != nil {
		os.Exit(2)
	}

	config := ReadConfig(*args.ConfigPath)
	if !ValidateConfig(config) {
		os.Exit(2)
	}

	mux := http.NewServeMux()
	s := Server{Config: config, IsDryRun: *args.IsDryRun}
	mux.Handle("/", s)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		slog.Error(err.Error())
		os.Exit(2)
	}
}
