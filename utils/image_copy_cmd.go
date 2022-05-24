//go:build cmd_implementation
// +build cmd_implementation

package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func CopyImage(image, registryURL string) (downloadFolder string, err error) {
	downloadFolder, err = ioutil.TempDir("", "temp-image*")
	if err != nil {
		return "", err
	}

	cmd, err := prepareCopyCmd(image, registryURL, downloadFolder)
	if err != nil {
		os.RemoveAll(downloadFolder)
		fmt.Println("Error while creating skopeo cmd\n Error:\n", err)

		return "", err
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		os.RemoveAll(downloadFolder)
		fmt.Println("Error when runing skopeo\n Output:\n", string(output))

		return "", err
	}

	return downloadFolder, nil
}

func prepareCopyCmd(image, registryAddr, downloadFolder string) (*exec.Cmd, error) {
	execPath, err := exec.LookPath("skopeo")
	if err != nil {
		return nil, err
	}

	imageTag := strings.Split(image, ":")[1]
	imageUrl := fmt.Sprintf("docker://%v/%v", registryAddr, image)
	ociLayout := fmt.Sprintf("oci:%v:%v", downloadFolder, imageTag)

	args := []string{
		"copy",
		"-f", "oci",
		"--src-tls-verify=false", // verify this from config
		imageUrl,                 // from
		ociLayout,                // to
	}

	cmd := exec.Command(execPath, args...)

	return cmd, nil
}
