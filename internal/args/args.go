package args

import (
	"flag"
)

// important to use pointers here, otherwise args will be set to default values
type Args struct {
	LogLevel   *string
	ConfigPath *string
	IsDryRun   *bool
}

func Parse() (a Args) {
	a.LogLevel = flag.String("log-level", "info", "debug | info | warn | error")
	a.ConfigPath = flag.String("config-path", "/app/config/config.yaml", "path to config file")
	a.IsDryRun = flag.Bool("dry-run", false, "will not execute commands")

	flag.Parse()
	return
}
