package action

import (
	"context"
	"errors"
	"github.com/hamster-shared/a-line-cli/pkg/model"
	"os"
	"os/exec"
	"strings"

	"github.com/hamster-shared/a-line-cli/pkg/logger"
	"github.com/hamster-shared/a-line-cli/pkg/output"
)

const STACK = "stack"

type DockerEnv struct {
	ctx         context.Context
	Image       string
	containerID string
	output      *output.Output
}

func NewDockerEnv(image string, ctx context.Context, output *output.Output) *DockerEnv {
	return &DockerEnv{
		ctx:    ctx,
		Image:  image,
		output: output,
	}
}

func (e *DockerEnv) Pre() error {
	e.output.NewStage("docker env")

	stack := e.ctx.Value(STACK).(map[string]interface{})

	data, ok := stack["workdir"]

	var workdir string
	if ok {
		workdir = data.(string)
	} else {
		return errors.New("workdir error")
	}

	workdirTmp := workdir + "@tmp"

	_ = os.MkdirAll(workdirTmp, os.ModePerm)

	//user := fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid())
	// "-u", user,

	commands := []string{"docker", "run", "-t", "-d", "-v", workdir + ":" + workdir, "-v", workdirTmp + ":" + workdirTmp, "-w", workdir, e.Image, "cat"}
	logger.Debugf("execute docker command: %s", strings.Join(commands, " "))
	e.output.WriteCommandLine(strings.Join(commands, " "))
	c := exec.Command(commands[0], commands[1:]...)
	output, err := c.CombinedOutput()
	if err != nil {
		logger.Errorf("execute docker command error: %s", err.Error())
		return err
	}
	containerID := string(output)
	logger.Debugf("docker command output: %s", containerID)
	e.output.WriteLine(containerID)

	e.containerID = strings.Fields(containerID)[0]
	return err
}

func (e *DockerEnv) Hook() (*model.ActionResult, error) {

	c := exec.Command("docker", "top", e.containerID, "-eo", "pid,comm")
	logger.Debugf("execute docker command: %s", strings.Join(c.Args, " "))
	e.output.WriteCommandLine(strings.Join(c.Args, " "))

	output, err := c.CombinedOutput()
	logger.Debugf("docker command output: %s", string(output))
	e.output.WriteLine(string(output))

	if err != nil {
		logger.Errorf("execute docker command error: %s", err.Error())
		return nil, err
	}

	stack := e.ctx.Value(STACK).(map[string]interface{})
	stack["withEnv"] = []string{"docker", "exec", e.containerID}
	return nil, nil
}

func (e *DockerEnv) Post() error {

	c := exec.Command("docker", "stop", "--time=1", e.containerID)
	logger.Debugf("execute docker command: %s", strings.Join(c.Args, " "))
	e.output.WriteCommandLine(strings.Join(c.Args, " "))

	output, err := c.CombinedOutput()
	e.output.WriteLine(string(output))

	if err != nil {
		logger.Errorf("execute docker command error: %s", err.Error())
		return err
	}

	c = exec.Command("docker", "rm", "-f", e.containerID)
	logger.Debugf("execute docker command: %s", strings.Join(c.Args, " "))
	e.output.WriteCommandLine(strings.Join(c.Args, " "))

	output, err = c.CombinedOutput()
	e.output.WriteLine(string(output))

	if err != nil {
		logger.Errorf("execute docker command error: %s", err.Error())
		return err
	}

	stack := e.ctx.Value(STACK).(map[string]interface{})
	stack["withEnv"] = []string{}

	return nil
}
