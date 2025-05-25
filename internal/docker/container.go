package docker

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/testcontainers/testcontainers-go/modules/compose"
)

func CreateStack(ctx context.Context , composeFile string) (*compose.DockerCompose, error) {
	log.Printf("Spinning up containers using compose file: %s\n", composeFile)
	// check if the compose file exists
	if !FileExists(composeFile) {
		return nil ,fmt.Errorf("Compose file does not exist: %v" , composeFile)
	}
	return compose.NewDockerComposeWith(compose.WithStackFiles(composeFile))
}

func SpinUpContainers(ctx context.Context,stack *compose.DockerCompose) error  {
	log.Println("Spinning up containers")
	return stack.Up(ctx)
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
