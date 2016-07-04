package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/containers/image/directory"
	"github.com/containers/image/docker"
	"github.com/containers/image/image"
	"github.com/containers/image/oci"
	"github.com/containers/image/openshift"
	"github.com/containers/image/types"
	"github.com/urfave/cli"
)

const (
	// atomicPrefix is the URL-like schema prefix used for Atomic registry image references.
	atomicPrefix = "atomic:"
	// dockerPrefix is the URL-like schema prefix used for Docker image references.
	dockerPrefix = "docker://"
	// directoryPrefix is the URL-like schema prefix used for local directories (for debugging)
	directoryPrefix = "dir:"
	// ociPrefix is the URL-like schema prefix used for OCI images.
	ociPrefix = "oci:"
)

// ParseImage converts image URL-like string to an initialized handler for that image.
func parseImage(c *cli.Context) (types.Image, error) {
	var (
		imgName   = c.Args().First()
		certPath  = c.GlobalString("cert-path")
		tlsVerify = c.GlobalBool("tls-verify")
	)
	switch {
	case strings.HasPrefix(imgName, dockerPrefix):
		return docker.NewImage(strings.TrimPrefix(imgName, dockerPrefix), certPath, tlsVerify)
		//case strings.HasPrefix(img, appcPrefix):
		//
	case strings.HasPrefix(imgName, directoryPrefix):
		src := directory.NewImageSource(strings.TrimPrefix(imgName, directoryPrefix))
		return image.FromSource(src, nil), nil
	}
	return nil, errors.New("no valid prefix provided")
}

// parseImageSource converts image URL-like string to an ImageSource.
func parseImageSource(c *cli.Context, name string) (types.ImageSource, error) {
	var (
		certPath  = c.GlobalString("cert-path")
		tlsVerify = c.GlobalBool("tls-verify") // FIXME!! defaults to false?
	)
	switch {
	case strings.HasPrefix(name, dockerPrefix):
		return docker.NewImageSource(strings.TrimPrefix(name, dockerPrefix), certPath, tlsVerify)
	case strings.HasPrefix(name, atomicPrefix):
		return openshift.NewImageSource(strings.TrimPrefix(name, atomicPrefix), certPath, tlsVerify)
	case strings.HasPrefix(name, directoryPrefix):
		return directory.NewImageSource(strings.TrimPrefix(name, directoryPrefix)), nil
	}
	return nil, fmt.Errorf("Unrecognized image reference %s", name)
}

// parseImageDestination converts image URL-like string to an ImageDestination.
func parseImageDestination(c *cli.Context, name string) (types.ImageDestination, error) {
	var (
		certPath  = c.GlobalString("cert-path")
		tlsVerify = c.GlobalBool("tls-verify") // FIXME!! defaults to false?
	)
	switch {
	case strings.HasPrefix(name, dockerPrefix):
		return docker.NewImageDestination(strings.TrimPrefix(name, dockerPrefix), certPath, tlsVerify)
	case strings.HasPrefix(name, atomicPrefix):
		return openshift.NewImageDestination(strings.TrimPrefix(name, atomicPrefix), certPath, tlsVerify)
	case strings.HasPrefix(name, directoryPrefix):
		return directory.NewImageDestination(strings.TrimPrefix(name, directoryPrefix)), nil
	case strings.HasPrefix(name, ociPrefix):
		return oci.NewImageDestination(strings.TrimPrefix(name, ociPrefix))
	}
	return nil, fmt.Errorf("Unrecognized image reference %s", name)
}
