package utils

import (
	"fmt"
	"os/exec"

	"zotregistry.io/zot/pkg/plugins/scan"
)

func ScanImage(imageName, imagePath string) (*scan.ScanResponse, error) {
	cmd, err := prepareScanCmd(imagePath)
	if err != nil {
		fmt.Println("Can't create scan cmd obj.", err)

		return emptyResponse(imageName), err
	}

	stdout, err := cmd.CombinedOutput()

	return ConvertScanOutputToRPCResponse(imageName, string(stdout)), err
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

func emptyResponse(image string) *scan.ScanResponse {
	return &scan.ScanResponse{
		Report: &scan.ScanReport{
			Image: &scan.Image{
				Name: image,
			},
			Scanner: &scan.Scanner{
				Name:    "ClamAv",
				Vendor:  "Cisco",
				Version: "0.1",
			},
			Vulnerabilities: make([]*scan.ScanVulnerability, 0),
		},
	}
}
