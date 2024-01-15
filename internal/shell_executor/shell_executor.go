package shell_executor

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
)

func Run(ctx context.Context, dir string, cmdline ...string) ([]byte, error) {
	var stderr, stdout bytes.Buffer

	cmd := exec.CommandContext(ctx, cmdline[0], cmdline[1:]...)
	cmd.Dir = dir
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	err := cmd.Run()

	if err != nil {
		return nil, fmt.Errorf("%s in %s cmd_error: %s err: %w", cmdline, dir, stderr.String(), err)
	}

	return stdout.Bytes(), nil
}
