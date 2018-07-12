package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gofn/gofn"
	"github.com/gofn/gofn/iaas/digitalocean"
	"github.com/gofn/gofn/provision"
)

func main() {
	contextDir := flag.String("contextDir", "./", "a string")
	dockerfile := flag.String("dockerfile", "Dockerfile", "a string")
	imageName := flag.String("imageName", "", "a string")
	remoteBuildURI := flag.String("remoteBuildURI", "", "a string")
	volumeSource := flag.String("volumeSource", "", "a string")
	volumeDestination := flag.String("volumeDestination", "", "a string")
	remoteBuild := flag.Bool("remoteBuild", false, "true or false")
	input := flag.String("input", "", "a string")
	flag.Parse()
	stdout, err := run(*contextDir, *dockerfile, *imageName, *remoteBuildURI, *volumeSource, *volumeDestination, *remoteBuild, *input)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	fmt.Println("Stdout: ", stdout)
}

func run(contextDir, dockerfile, imageName, remoteBuildURI, volumeSource, volumeDestination string, remote bool, input string) (stdout string, err error) {
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
	if remote {
		key := os.Getenv("DIGITALOCEAN_API_KEY")
		if key == "" {
			log.Fatalln("You must provide an api key for digital ocean")
		}
		do, err := digitalocean.New(key)
		if err != nil {
			log.Println(err)
		}
		buildOpts.Iaas = do
	}
	stdout, stderr, err := gofn.Run(context.Background(), buildOpts, containerOpts)
	if stderr != "" {
		stdout = stderr
	}
	return
}
