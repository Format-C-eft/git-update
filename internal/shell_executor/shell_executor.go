package shell_executor

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"regexp"
)

var regexpLineBreaks = regexp.MustCompile(`\r?\n`)

func Run(ctx context.Context, dir string, cmdline ...string) (string, error) {
	var stderr, stdout bytes.Buffer

	cmd := exec.CommandContext(ctx, cmdline[0], cmdline[1:]...)
	cmd.Dir = dir
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	err := cmd.Start()
	if err != nil {
		return "", fmt.Errorf("%s in %s cmd_error: %s err: %w",
			cmdline,
			dir,
			string(regexpLineBreaks.ReplaceAll(stderr.Bytes(), []byte(" "))),
			err,
		)
	}

	err = cmd.Wait()
	if err != nil {
		return "", fmt.Errorf("%s in %s cmd_error: %s err: %w",
			cmdline,
			dir,
			string(regexpLineBreaks.ReplaceAll(stderr.Bytes(), []byte(" "))),
			err,
		)
	}

	return string(regexpLineBreaks.ReplaceAll(stdout.Bytes(), []byte(" "))), nil
}
