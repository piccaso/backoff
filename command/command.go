package command

import (
	"context"
	"errors"
	"os"
	"os/exec"
)

func NewCommandWithContext(ctx context.Context, args []string) (*exec.Cmd, context.Context, error) {
	argsLen := len(args)

	if argsLen < 1 {
		return nil, ctx, errors.New("args missing")
	}

	if ctx == nil {
		return nil, ctx, errors.New("ctx is nil")
	}

	var cmd *exec.Cmd
	if argsLen == 1 {
		cmd = exec.CommandContext(ctx, args[0])
	} else {
		cmd = exec.CommandContext(ctx, args[0], args[1:]...)
	}

	wirePipes(cmd)
	return cmd, ctx, nil
}

func wirePipes(cmd *exec.Cmd) {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
}
