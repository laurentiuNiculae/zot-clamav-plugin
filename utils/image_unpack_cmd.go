//go:build cmd_implementation
// +build cmd_implementation

package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

func UnpackImage(imagePath string) (unpackDir string, err error) {
	unpackDir, err = ioutil.TempDir("", "zot-clamav-scanner-unpack-dir*")
	if err != nil {
		defer os.RemoveAll(unpackDir)

		return "", err
	}

	cmd, err := prepareUnpackCmd(imagePath, unpackDir)
	if err != nil {
		defer os.RemoveAll(unpackDir)

		return "", err
	}

	b, err := cmd.CombinedOutput()
	if err != nil {
		defer os.RemoveAll(unpackDir)
		fmt.Printf("Output:\n %v", string(b))

		return "", err
	}

	return unpackDir, nil
}

func prepareUnpackCmd(imagePath, unpackDir string) (*exec.Cmd, error) {
	execPath, err := exec.LookPath("umoci")
	if err != nil {
		return nil, err
	}

	args := []string{
		"unpack",
		"--rootless",
		"--image", imagePath,
		unpackDir,
	}

	cmd := exec.Command(execPath, args...)

	return cmd, nil
}
