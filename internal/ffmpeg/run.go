package ffmpeg

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// Runner defines how external commands are executed.
type Runner interface {
	Run(ctx context.Context, name string, args []string) (stdout, stderr []byte, err error)
}

// ExecRunner executes processes via os/exec without shell.
type ExecRunner struct{}

func (ExecRunner) Run(ctx context.Context, name string, args []string) ([]byte, []byte, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err := cmd.Run()
	return outBuf.Bytes(), errBuf.Bytes(), err
}

// LookPath verifies a binary is present.
func LookPath(bin string) (string, error) {
	p, err := exec.LookPath(bin)
	if err != nil {
		return "", fmt.Errorf("%s not found in PATH", bin)
	}
	return p, nil
}

// PrettyCmd renders a friendly representation of a command for errors/logs.
func PrettyCmd(name string, args []string) string {
	var b strings.Builder
	b.WriteString(name)
	for _, a := range args {
		b.WriteByte(' ')
		// naive quoting for display only (not for execution)
		if strings.ContainsAny(a, " \t\n'\"") {
			b.WriteByte('"')
			b.WriteString(strings.ReplaceAll(a, "\"", "\\\""))
			b.WriteByte('"')
		} else {
			b.WriteString(a)
		}
	}
	return b.String()
}
