package action

import (
	"bufio"
	"context"
	"errors"
	"os"
	"os/exec"
	"state-example/pkg/logger"
	"state-example/pkg/stream"
	"state-example/pkg/utils"
	"strings"
	"syscall"
)

// ShellAction 命令工作
type ShellAction struct {
	command  string
	filename string
	ctx      context.Context
}

func NewShellAction(command string, ctx context.Context) *ShellAction {

	return &ShellAction{
		command: command,
		ctx:     ctx,
	}
}

func (a *ShellAction) Pre() error {

	stack := a.ctx.Value(STACK).(map[string]interface{})

	data, ok := stack["workdir"]

	var workdir string
	if ok {
		workdir = data.(string)
	} else {
		return errors.New("workdir error")
	}

	workdirTmp := workdir + "@tmp"

	_ = os.MkdirAll(workdirTmp, os.ModePerm)

	a.filename = workdirTmp + "/" + utils.RandSeq(10) + ".sh"

	content := []byte("#!/bin/sh\nset -ex\n" + a.command)
	err := os.WriteFile(a.filename, content, os.ModePerm)

	return err
}

func (a *ShellAction) Hook() error {

	stack := a.ctx.Value(STACK).(map[string]interface{})

	workdir, ok := stack["workdir"].(string)
	if !ok {
		return errors.New("workdir is empty")
	}
	logger.Infof("shell stack: %v", stack)

	commands := []string{"sh", "-c", a.filename}
	val, ok := stack["withEnv"]
	if ok {
		precommand := val.([]string)
		shellCommand := make([]string, len(commands))
		copy(shellCommand, commands)
		commands = append([]string{}, precommand...)
		commands = append(commands, shellCommand...)
	}

	c := exec.CommandContext(a.ctx, commands[0], commands[1:]...) // mac linux
	c.Dir = workdir
	logger.Debugf("execute shell command: %s", strings.Join(commands, " "))

	stdout, err := c.StdoutPipe()
	if err != nil {
		logger.Errorf("stdout error: %v", err)
		return err
	}
	stderr, err := c.StderrPipe()
	if err != nil {
		logger.Errorf("stderr error: %v", err)
		return err
	}

	go func() {
		for {
			// 检测到 ctx.Done() 之后停止读取
			<-a.ctx.Done()
			if a.ctx.Err() != nil {
				logger.Errorf("shell command error: %v", a.ctx.Err())
			} else {
				p := c.Process
				if p == nil {
					return
				}
				// Kill by negative PID to kill the process group, which includes
				// the top-level process we spawned as well as any subprocesses
				// it spawned.
				_ = syscall.Kill(-p.Pid, syscall.SIGKILL)
				logger.Info("shell command killed")
			}
		}
	}()

	stdoutScanner := bufio.NewScanner(stdout)
	stderrScanner := bufio.NewScanner(stderr)
	go func() {
		for stdoutScanner.Scan() {
			stream.OutputCh <- stdoutScanner.Text()
		}
	}()
	go func() {
		for stderrScanner.Scan() {
			stream.OutputCh <- stderrScanner.Text()
		}
	}()

	err = c.Start()
	if err != nil {
		logger.Errorf("shell command start error: %v", err)
		return err
	}

	err = c.Wait()
	if err != nil {
		logger.Errorf("shell command wait error: %v", err)
		return err
	}

	logger.Info("execute shell command success")
	return err
}

func (a *ShellAction) Post() error {
	return os.Remove(a.command)
}
