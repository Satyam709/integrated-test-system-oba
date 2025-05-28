package static_tests

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/Satyam709/integrated-test-system-oba/internal/docker"
)

var (
	DockerManager *docker.DockerManager
	ctx           context.Context
)

func TestMain(m *testing.M) {
	// Set up the environment for static tests
	// This could include setting up environment variables, initializing configurations, etc.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	composeFile := filepath.Join(os.Getenv("PROJECT_ROOT"), "docker", "docker-compose.yml")
	// Call the function to spin up containers
	DockerManager = docker.GetInstance()
	err := DockerManager.InitializeStack(ctx, composeFile)
	if err!= nil {
		log.Fatalf("Error initializing stack: %v", err)
	}
	stack := DockerManager.GetStack()

	obaContainer,err := DockerManager.GetOBAServerContainer(ctx)
	if err != nil {
		log.Fatalf("Error getting OBA server container: %v", err)
	}

	// build_bundle.sh log
	err = obaContainer.BuildBundle(ctx)
	if err != nil {
		log.Fatalf("Error building GTFS bundle: %v", err)
	}
	//start the OBA server
	err = obaContainer.StartServer(ctx)
	if err != nil {
		log.Fatalf("Error starting OBA server: %v", err)
	}
	// Wait for port 8080 (inside container)
	err = obaContainer.WaitForServerReady(ctx)
	if err != nil {
		log.Fatalf("Error waiting for OBA server to be ready: %v", err)
	}
	// Run the tests
	fmt.Printf("Running static tests...\n")
	exitCode := m.Run()
	fmt.Printf("exitCode: %v\n", exitCode)
	// Clean up the environment after tests
	err = stack.Down(ctx)
	if err != nil {
		fmt.Printf("Error tearing down containers: %v\n", err)
	}
}
