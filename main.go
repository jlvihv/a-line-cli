package main

import (
	"github.com/hamster-shared/a-line-cli/cmd"
	"github.com/hamster-shared/a-line-cli/pkg/logger"
	"github.com/hamster-shared/a-line-cli/pkg/stream"

	"github.com/sirupsen/logrus"
)

func main() {
	logger.Init().ToStdout().SetLevel(logrus.TraceLevel)
	go stream.Output()
	cmd.Execute()
}
