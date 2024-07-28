package main

import (
	"actioneer/internal/args"
	"actioneer/internal/command"
	"actioneer/internal/config"
	"actioneer/internal/logging"
	"actioneer/internal/notification"
	"actioneer/internal/processor"
	"actioneer/internal/state"
	"io"
	"log/slog"
	"net/http"
	"os"
)

type Server struct {
	IsDryRun bool
	State    state.State
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bytes, errReadBody := io.ReadAll(r.Body)
	if errReadBody != nil {
		panic(errReadBody)
	}
	defer r.Body.Close()

	notificationExternal, errReadExternalNotification := notification.ReadAlertmanagerNotification(bytes)
	if errReadExternalNotification != nil {
		panic(errReadExternalNotification)
	}

	notification := notification.ToInternal(notificationExternal, s.State)
	
	shell := command.CommandRunner{}
	errTakeAction := processor.TakeActions(shell, s.State, notification, s.IsDryRun)
	if errTakeAction != nil {
		panic(errTakeAction)
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
