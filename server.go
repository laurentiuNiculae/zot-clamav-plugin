package main

import (
	"context"
	"fmt"
	"os"

	"zotregistry.io/zot/pkg/plugins/scan"
)

type ScanServer struct {
	scan.UnimplementedScanServer
}

func (ss ScanServer) Scan(ctx context.Context, request *scan.ScanRequest) (*scan.ScanResponse, error) {
	var (
		image       = request.GetImage()
		registryURL = request.GetRegistry().GetUrl()
	)

	// download the image using skopeo
	downloadDir, err := CopyImage(image, registryURL)
	if err != nil {
		return emptyReport(request.Image), nil
	}
	defer os.RemoveAll(downloadDir)

	// unpack the image with umoci
	unpackDir, err := UnpackImage(downloadDir)
	if err != nil {
		fmt.Println("Error when unpacking", err)

		return emptyReport(request.Image), nil
	}
	defer os.RemoveAll(unpackDir)

	// scan using clamav
	scanResult, err := ScanImage(unpackDir)
	if err != nil {
		fmt.Println("Error when scanning", err)

		return emptyReport(request.Image), nil
	}

	fmt.Println(scanResult)

	return emptyReport(request.Image), nil
}

func emptyReport(image string) *scan.ScanResponse {
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
