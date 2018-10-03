package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/gofn/gofn"
	"github.com/gofn/gofn/provision"
)

func main() {
	contextDir := flag.String("contextDir", "./", "a string")
	dockerfile := flag.String("dockerfile", "Dockerfile", "a string")
	imageName := flag.String("imageName", "", "a string")
	remoteBuildURI := flag.String("remoteBuildURI", "", "a string")
	volumeSource := flag.String("volumeSource", "", "a string")
	volumeDestination := flag.String("volumeDestination", "", "a string")
	input := flag.String("input", "", "a string")
	flag.Parse()
	stdout, err := run(*contextDir, *dockerfile, *imageName, *remoteBuildURI, *volumeSource, *volumeDestination, *input)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	fmt.Println("Stdout: ", stdout)
}

func run(contextDir, dockerfile, imageName, remoteBuildURI, volumeSource, volumeDestination string, input string) (stdout string, err error) {
	buildOpts := &provision.BuildOptions{
		ContextDir: contextDir,
		Dockerfile: dockerfile,
		ImageName:  imageName,
		RemoteURI:  remoteBuildURI,
		StdIN:      input,
	}
	containerOpts := &provision.ContainerOptions{}
	if volumeSource != "" {
		if volumeDestination == "" {
			volumeDestination = volumeSource
		}
		containerOpts.Volumes = []string{fmt.Sprintf("%s:%s", volumeSource, volumeDestination)}
	}
	// TODO: remote (IaaS) support
	stdout, stderr, err := gofn.Run(context.Background(), buildOpts, containerOpts)
	if stderr != "" {
		stdout = stderr
	}
	return
}
