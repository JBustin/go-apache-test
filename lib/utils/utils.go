package utils

import (
	"bytes"
	"os/exec"
)

func Exec(name string, args ...string) (string, string, error) {
	cmd := exec.Command(name, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	return stdout.String(), stderr.String(), err
}

func ExtendStr(args ...string) string {
	for _, v := range args {
		if v != "" {
			return v
		}
	}
	return args[0]
}

func ExtendInt(args ...int) int {
	for _, v := range args {
		if v != 0 {
			return v
		}
	}
	return args[0]
}
