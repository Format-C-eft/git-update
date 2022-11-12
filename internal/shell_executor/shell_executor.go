package shell_executor

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"

	"github.com/pkg/errors"
)

func Run(ctx context.Context, dir string, cmdline ...string) ([]byte, error) {
	var stderr, stdout bytes.Buffer

	cmd := exec.CommandContext(ctx, cmdline[0], cmdline[1:]...)
	cmd.Dir = dir
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	err := cmd.Run()

	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("%s in %s error: %s", cmdline, dir, stderr.String()))
	}

	return stdout.Bytes(), nil
}
