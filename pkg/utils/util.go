package utils

import (
	"context"
	"math/rand"
	"os/exec"
	"os/user"
	"syscall"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GetUIDAndGID() (string, string, error) {
	u, err := user.Current()
	if err != nil {
		return "", "", err
	}
	return u.Uid, u.Gid, nil
}

type Cmd struct {
	ctx context.Context
	*exec.Cmd
}

// NewCommand is like exec.CommandContext but ensures that subprocesses
// are killed when the context times out, not just the top level process.
func NewCommand(ctx context.Context, command string, args ...string) *Cmd {
	return &Cmd{ctx, exec.Command(command, args...)}
}

func (c *Cmd) Start() error {
	// Force-enable setpgid bit so that we can kill child processes when the
	// context times out or is canceled.
	if c.Cmd.SysProcAttr == nil {
		c.Cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	c.Cmd.SysProcAttr.Setpgid = true
	err := c.Cmd.Start()
	if err != nil {
		return err
	}
	go func() {
		<-c.ctx.Done()
		p := c.Cmd.Process
		if p == nil {
			return
		}
		// Kill by negative PID to kill the process group, which includes
		// the top-level process we spawned as well as any subprocesses
		// it spawned.
		_ = syscall.Kill(-p.Pid, syscall.SIGKILL)
	}()
	return nil
}

func (c *Cmd) Run() error {
	if err := c.Start(); err != nil {
		return err
	}
	return c.Wait()
}
