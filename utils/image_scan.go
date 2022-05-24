//go:build !cmd_implementation
// +build !cmd_implementation

package utils

import (
	"fmt"
	"os/exec"
)

func ScanImage(imagePath string) (string, error) {
	cmd, err := prepareScanCmd(imagePath)
	if err != nil {
		fmt.Println("Can't create scan cmd obj.", err)

		return "", err
	}

	stdout, err := cmd.CombinedOutput()

	return string(stdout), err
}

func prepareScanCmd(imagePath string) (*exec.Cmd, error) {
	execPath, err := exec.LookPath("clamdscan")
	if err != nil {
		return nil, err
	}

	args := []string{
		"--fdpass",
		"-m",
		imagePath,
	}

	cmd := exec.Command(execPath, args...)

	return cmd, nil
}
