package server

import (
	"context"
	"os"

	"github.com/laurentiuNiculae/zot-clamav-plugin/utils"
	"zotregistry.io/zot/pkg/plugins/scan"
)

type ScanServer struct {
	scan.UnimplementedScanServer
}

func (ss ScanServer) Scan(ctx context.Context, request *scan.ScanRequest) (*scan.ScanResponse, error) {
	var (
		imageName   = request.GetImage()
		registryURL = request.GetRegistry().GetUrl()
	)

	downloadDir, err := utils.CopyImage(imageName, registryURL)
	if err != nil {
		return emptyRespose(request.Image), nil
	}
	defer os.RemoveAll(downloadDir)

	// unpack the image with umoci
	// unpackDir, err := utils.UnpackImage(image, downloadDir)
	// if err != nil {
	// 	fmt.Println("Error when unpacking: ", err)

	// 	return emptyReport(request.Image), nil
	// }
	// defer os.RemoveAll(unpackDir)

	// scan using clamav
	scanResponse, err := utils.ScanImage(imageName, downloadDir)

	return scanResponse, nil
}

func emptyRespose(image string) *scan.ScanResponse {
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
