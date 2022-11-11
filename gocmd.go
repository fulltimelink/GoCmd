package GoCmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type goCmd struct {
	cmd  *exec.Cmd
	args []string
}

type Option func(*goCmd)

var CommandNotFoundErr = errors.New("command not found")

func WithStdin(stdin io.Reader) Option {
	return func(c *goCmd) {
		c.cmd.Stdin = stdin
	}
}

func WithStdout(stdout io.Writer) Option {
	return func(c *goCmd) {
		c.cmd.Stdout = stdout
	}
}

func WithStderr(stderr io.Writer) Option {
	return func(c *goCmd) {
		c.cmd.Stderr = stderr
	}
}

func WithOutErr(outerr io.Writer) Option {
	return func(c *goCmd) {
		c.cmd.Stdout = outerr
		c.cmd.Stderr = outerr
	}
}

func WithEnv(key, value string) Option {
	return func(c *goCmd) {
		c.cmd.Env = append(os.Environ(), fmt.Sprintf("%s=%s", key, value))
	}
}

func WithArgs(args []string) Option {
	return func(c *goCmd) {
		c.args = args
	}
}

func applyOptions(cmd *goCmd, opts []Option) {
	for _, opt := range opts {
		opt(cmd)
	}
}

func initGoCmd(cmd string, opts ...Option) (*exec.Cmd, error) {
	_, err := exec.LookPath(cmd)
	if nil != err {
		return nil, CommandNotFoundErr
	}
	command := exec.Command(cmd)

	gcmd := goCmd{
		cmd:  command,
		args: nil,
	}
	applyOptions(&gcmd, opts)
	if nil != gcmd.args {
		gcmd.cmd = exec.Command(cmd, gcmd.args...)
	}
	return gcmd.cmd, nil
}

// 执行命令
func RunCommand(cmd string, opts ...Option) error {
	command, err := initGoCmd(cmd, opts...)
	if nil != err {
		return err
	}
	return command.Run()
}
