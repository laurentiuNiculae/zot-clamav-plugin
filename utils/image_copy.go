//go:build !cmd_implementation
// +build !cmd_implementation

package utils

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/containers/image/v5/copy"
	"github.com/containers/image/v5/signature"
	"github.com/containers/image/v5/transports/alltransports"
	"github.com/containers/image/v5/types"
	imgspecv1 "github.com/opencontainers/image-spec/specs-go/v1"
)

func CopyImage(image, registryURL string) (downloadFolder string, err error) {
	downloadFolder, err = ioutil.TempDir("", "temp-image*")
	if err != nil {
		return "", err
	}

	ctx := context.Background() // maby something with timeout

	imageTag := strings.Split(image, ":")[1]
	imageUrl := fmt.Sprintf("docker://%v/%v", registryURL, image)
	ociLayout := fmt.Sprintf("oci:%v:%v", downloadFolder, imageTag)

	srcRef, err := alltransports.ParseImageName(imageUrl)
	if err != nil {
		return "", fmt.Errorf("invalid source name %s: %v", imageUrl, err)
	}

	destRef, err := alltransports.ParseImageName(ociLayout)
	if err != nil {
		return "", fmt.Errorf("invalid destination name %s: %v", ociLayout, err)
	}

	sourceCtx := getSourceCtx()

	destinationCtx := getDestCtx()

	policyContext, err := signature.NewPolicyContext(
		&signature.Policy{
			Default: []signature.PolicyRequirement{
				signature.NewPRInsecureAcceptAnything(),
			},
		},
	)
	if err != nil {
		fmt.Println("Error when creating now policy context", err)

		return "", err
	}

	manifestBytes, err := copy.Image(ctx, policyContext, destRef, srcRef, &copy.Options{
		SourceCtx:             sourceCtx,
		DestinationCtx:        destinationCtx,
		ForceManifestMIMEType: imgspecv1.MediaTypeImageManifest,
		ImageListSelection:    copy.CopySystemImage,
	})
	if err != nil {
		return "", err
	}

	fmt.Println(string(manifestBytes))

	return downloadFolder, nil
}

func getDestCtx() *types.SystemContext {
	ctx := &types.SystemContext{}
	ctx.DockerInsecureSkipTLSVerify = types.NewOptionalBool(true)

	return ctx
}

func getSourceCtx() *types.SystemContext {
	ctx := &types.SystemContext{}
	ctx.DockerInsecureSkipTLSVerify = types.NewOptionalBool(true)

	return ctx
}
