package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
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

	var notification Notification
	if err := json.Unmarshal(bytes, &notification); err != nil {
		slog.Error("cannot unmarshal incomming notification: " + fmt.Sprint(notification))
		slog.Error(err.Error())
	}

	slog.Debug("processing notification: " + fmt.Sprint(notification))
	if len(notification.Alerts) > 0 {
		for _, alert := range notification.Alerts {
			if alertName, ok := alert.Labels[alertNameLabel]; ok {
				slog.Debug("processing alert [" + alertName + "]")
				for _, action := range s.Config.Actions {
					if action.Alertname == alertName {
						slog.Debug("processing action: [" + fmt.Sprint(action) + "]")
						if s.IsDryRun {
							slog.Info("dry run: [" + action.Command + "]")
						} else {
							// if err := execCommand(action.Command); err != nil {
							// 	slog.Error("command execution failed: [" + action.Command + "]")
							// 	slog.Error(err.Error())
							// }
						}
					}
				}
			} else {
				slog.Error("no alert name label [" + alertNameLabel + "], skipping: " + fmt.Sprintf("%+v", alert))
			}
		}
	}
}

func InitLogger(log_level string) {
	var programLogLevel = new(slog.LevelVar)
	h := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: programLogLevel})
	slog.SetDefault(slog.New(h))

	switch log_level {
	case "debug":
		programLogLevel.Set(slog.LevelDebug)
	case "info":
		programLogLevel.Set(slog.LevelInfo)
	case "warn":
		programLogLevel.Set(slog.LevelWarn)
	case "error":
		programLogLevel.Set(slog.LevelError)
	default:
		slog.Error("wrong value in --log-level=" + log_level)
		os.Exit(2)
	}
	slog.Info("--log-level=" + strings.ToLower(programLogLevel.Level().String()))
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
	log_level := flag.String("log-level", "info", "debug | info | warn | error")
	config_path := flag.String("config-path", "/app/config/config.json", "path to config file")
	is_dry_run := flag.Bool("dry-run", false, "will not execute commands")
	flag.Parse()

	InitLogger(*log_level)

	config := ReadConfig(*config_path)
	if !ValidateConfig(config) {
		os.Exit(2)
	}

	mux := http.NewServeMux()
	s := Server{Config: config, IsDryRun: *is_dry_run}
	mux.Handle("/", s)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		slog.Error(err.Error())
		os.Exit(2)
	}
}
