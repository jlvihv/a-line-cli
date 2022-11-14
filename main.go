package main

import (
	"state-example/cmd"
	"state-example/pkg/logger"

	"github.com/sirupsen/logrus"
)

func main() {
	logger.Init().ToStdout().SetLevel(logrus.TraceLevel)
	cmd.Execute()
}
