package utils

import (
	"context"
	"errors"
	"os/exec"
	"time"
)

// ExecuteCommand ...
func ExecuteCommand(command string, arguments []string, dir string, env []string, timeout time.Duration) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel() // The cancel should be deferred so resources are cleaned up

	cmd := exec.CommandContext(ctx, command, arguments...)
	cmd.Dir = dir
	cmd.Env = env
	out, err := cmd.CombinedOutput()

	if ctx.Err() == context.DeadlineExceeded {
		return out, errors.New(`Init Timeout was reached.`)
	}

	return out, err
}
