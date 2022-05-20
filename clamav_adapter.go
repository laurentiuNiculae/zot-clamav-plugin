package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
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
		imagePath,
	}

	cmd := exec.Command(execPath, args...)

	return cmd, nil
}

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
