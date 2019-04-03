package main

import (
	"os"

	"github.com/mackerelio/golib/logging"

	"github.com/a-know/mackerel-remora/agent"
)

const cmdName = "mackerel-remora"

var version, revision string

var logger = logging.GetLogger("main")

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	logger.Infof("starting %s (version:%s, revision:%s)", cmdName, version, revision)
	if err := agent.NewAgent(version, revision).Run(args); err != nil {
		logger.Errorf("%s", err)
		return 1
	}
	return 0
}
