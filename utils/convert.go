package utils

import (
	"regexp"
	"strings"

	"zotregistry.io/zot/pkg/plugins/scan"
)

func ConvertScanOutputToRPCResponse(imageName, stdOut string) *scan.ScanResponse {
	return &scan.ScanResponse{
		Report: &scan.ScanReport{
			Image: &scan.Image{
				Name: imageName,
			},
			Scanner: &scan.Scanner{
				Name: "ClamAv",
			},
			Vulnerabilities: ConvertToRPCVuln(stdOut),
		},
	}
}

func ConvertToRPCVuln(stdOut string) []*scan.ScanVulnerability {
	const (
		filePath  = 1
		virusName = 2
	)

	// example of vuln:
	// img/blobs/sha256/e228 ... ee36: Win.Test.EICAR_HDB-1 FOUND
	vulnsFinder, _ := regexp.Compile(`(.*): (.*) FOUND\n`)
	foundVulns := vulnsFinder.FindAllStringSubmatch(stdOut, -1)

	ScanVulnerabilities := make([]*scan.ScanVulnerability, len(foundVulns))

	for i, vuln := range foundVulns {
		ScanVulnerabilities[i] = &scan.ScanVulnerability{
			VulnerabilityId: "-",
			Pkg:             "",
			Version:         "",
			FixedVersion:    "",
			Title:           vuln[virusName],
			Description:     "",
			Severity:        scan.Severity_HIGH,
			Layer: &scan.Layer{
				Digest: ConvertFilePathToDigest(vuln[filePath]),
				DiffId: "",
			},
		}
	}

	return ScanVulnerabilities
}

func ConvertFilePathToDigest(s string) string {
	digestFinder, _ := regexp.Compile("blobs/(.*)/(.*)")

	submatch := digestFinder.FindStringSubmatch(s)

	if submatch == nil {
		return ""
	}

	return strings.Join(submatch[1:], "")
}
