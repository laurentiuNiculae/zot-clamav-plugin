//go:build !cmd_implementation
// +build !cmd_implementation

package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/opencontainers/umoci"
	"github.com/opencontainers/umoci/oci/cas/dir"
	"github.com/opencontainers/umoci/oci/casext"
	"github.com/opencontainers/umoci/oci/layer"
	"github.com/opencontainers/umoci/pkg/idtools"
	"github.com/pkg/errors"
)

func UnpackImage(imageName, imagePath string) (unpackDir string, err error) {
	unpackDir, err = ioutil.TempDir("", "zot-clamav-scanner-unpack-dir*")
	if err != nil {
		defer os.RemoveAll(unpackDir)

		return "", err
	}
	imageTag := strings.Split(imageName, ":")[1]

	var unpackOptions layer.UnpackOptions
	var meta umoci.Meta
	meta.Version = umoci.MetaVersion
	meta.MapOptions.Rootless = true

	uidmapidMap, err := idtools.ParseMapping(fmt.Sprintf("0:%d:1", os.Geteuid()))
	if err != nil {
		return "", errors.Wrapf(err, "failure parsing --uid-map %s", uidmapidMap)
	}

	gidmap, err := idtools.ParseMapping(fmt.Sprintf("0:%d:1", os.Getegid()))
	if err != nil {
		return "", errors.Wrapf(err, "failure parsing --gid-map %s", gidmap)
	}

	meta.MapOptions.UIDMappings = append(meta.MapOptions.UIDMappings, uidmapidMap)
	meta.MapOptions.GIDMappings = append(meta.MapOptions.GIDMappings, gidmap)

	unpackOptions.KeepDirlinks = false
	unpackOptions.MapOptions = meta.MapOptions

	// Get a reference to the CAS.
	engine, err := dir.Open(imagePath)
	if err != nil {
		return "", err
	}

	engineExt := casext.NewEngine(engine)
	defer engine.Close()

	err = umoci.Unpack(engineExt, imageTag, unpackDir, unpackOptions)
	if err != nil {
		return "", err
	}

	return unpackDir, nil
}
