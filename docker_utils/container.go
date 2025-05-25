package docker_utils

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/testcontainers/testcontainers-go/modules/compose"
)

// func SpinContainers(composeFile string) {
// 	// check if the compose file exists
// 	if !FileExists(composeFile) {
// 		panic("Compose file does not exist: " + composeFile)
// 	}
// 	// Spin up the containers using the docker-compose command
// 	cmd := "docker-compose -f " + composeFile + " up -d"
// 	output, err := RunCommand(cmd)
// 	if err != nil {
// 		panic("Failed to spin up containers: " + err.Error())
// 	}
// 	if output != "" {
// 		println("Output from spinning up containers: " + output)
// 	} else {
// 		println("Containers spun up successfully.")
// 	}
// }

func SpinUpContainers(ctx context.Context , composeFile string) *compose.DockerCompose {
	fmt.Printf("Spinning up containers using compose file: %s\n", composeFile)

	// check if the compose file exists
	if !FileExists(composeFile) {
		panic("Compose file does not exist: " + composeFile)
	}

	stack, err := compose.NewDockerComposeWith(compose.WithStackFiles(composeFile))
	if err != nil {
		panic("Failed to create Docker Compose stack: " + err.Error())
	}

	err = stack.Up(ctx)
	if err != nil {
		panic("Failed to spin up containers: " + err.Error())
	}
	return stack
}

func RunCommand(cmd string) (string, error) {
	// Execute the commands
	command := exec.Command("sh", "-c", cmd)
	output, err := command.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func FileExists(composeFile string) bool {
	// Check if the file exists
	_, err := os.Stat(composeFile)
	return !os.IsNotExist(err)
}
