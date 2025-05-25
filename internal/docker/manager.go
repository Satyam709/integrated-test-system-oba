package docker

import (
	"context"
	"log"
	"sync"

	"github.com/testcontainers/testcontainers-go/modules/compose"
)

// DockerManager manages Docker operations as a singleton
type DockerManager struct {
	stack *compose.DockerCompose
}

var (
	instance *DockerManager
	once     sync.Once
)

// GetInstance returns the singleton instance of DockerManager
func GetInstance() *DockerManager {
	once.Do(func() {
		instance = &DockerManager{}
	})
	return instance
}

// InitializeStack sets up the Docker Compose stack and spins up the containers
func (dm *DockerManager) InitializeStack(ctx context.Context, composeFile string) error {
	stack, err := CreateStack(ctx, composeFile)
	if err!= nil {
		log.Printf("failed to create stack")
		return err
	}
	err = SpinUpContainers(ctx, stack)
	if err!= nil {
		log.Printf("failed to spin up containers")
		return err
	}
	dm.stack = stack
	return nil
}

// GetStack returns the current Docker Compose stack
func (dm *DockerManager) GetStack() *compose.DockerCompose {
	return dm.stack
}

// Cleanup tears down the Docker Compose stack
func (dm *DockerManager) Cleanup(ctx context.Context) error {
	if dm.stack != nil {
		return dm.stack.Down(ctx)
	}
	return nil
}