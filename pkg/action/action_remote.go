package action

import (
	"context"
	"fmt"
	"github.com/hamster-shared/a-line-cli/pkg/logger"
	"github.com/hamster-shared/a-line-cli/pkg/model"
	"github.com/hamster-shared/a-line-cli/pkg/utils"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"os/exec"
	"path"
)

// RemoteAction 执行远程命令
type RemoteAction struct {
	name string
	args map[string]string
	ctx  context.Context

	actionRoot string
}

func NewRemoteAction(name string, args map[string]string, ctx context.Context) *RemoteAction {
	return &RemoteAction{
		name: name,
		args: args,
		ctx:  ctx,
	}
}

func (a *RemoteAction) Pre() error {

	stack := a.ctx.Value(STACK).(map[string]interface{})

	pipelineName := stack["name"].(string)

	logger.Infof("git stack: %v", stack)

	hamsterRoot := stack["hamsterRoot"].(string)
	cloneDir := utils.RandSeq(9)
	a.actionRoot = path.Join(hamsterRoot, cloneDir)

	_ = os.MkdirAll(hamsterRoot, os.ModePerm)
	_ = os.Remove(path.Join(hamsterRoot, pipelineName))

	githubUrl := fmt.Sprintf("https://github.com/%s", a.name)

	commands := []string{"git", "clone", "--progress", githubUrl, cloneDir}
	c := exec.CommandContext(a.ctx, commands[0], commands[1:]...) // mac linux
	c.Dir = hamsterRoot

	fmt.Println(a.name)
	fmt.Println(a.args)

	output, err := c.CombinedOutput()
	fmt.Println(string(output))
	return err
}

func (a *RemoteAction) Hook() error {

	file, err := os.Open(path.Join(a.actionRoot, "action.yml"))
	if err != nil {
		return err
	}
	yamlFile, err := io.ReadAll(file)
	var remoteAction model.RemoteAction
	err = yaml.Unmarshal(yamlFile, &remoteAction)

	for index, step := range remoteAction.Runs.Steps {
		args := make([]string, 0)
		if _, err := os.Stat(path.Join(a.actionRoot, step.Run)); err == nil {
			args = append(args, step.Run)
		} else {
			stepFile := path.Join(a.actionRoot, fmt.Sprintf("step-%d", index))
			err = os.WriteFile(stepFile, []byte(step.Run), os.ModePerm)
			args = append(args, "-c", stepFile)
		}

		cmd := utils.NewCommand(a.ctx, step.Shell, args...)
		cmd.Dir = a.actionRoot
		output, _ := cmd.CombinedOutput()
		fmt.Println(string(output))
	}

	return err
}

func (a *RemoteAction) Post() error {

	return os.RemoveAll(a.actionRoot)
}
