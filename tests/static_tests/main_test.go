package static_tests

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	onebusaway "github.com/OneBusAway/go-sdk"
	"github.com/OneBusAway/go-sdk/option"
	"github.com/Satyam709/integrated-test-system-oba/internal/docker"
)

var (
	DockerManager *docker.DockerManager
	ctx           context.Context
	obaClient     *onebusaway.Client
	obaContainer  *docker.OBAContainer
)

func TestMain(m *testing.M) {
	// Set up the environment for static tests

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	composeFile := filepath.Join(os.Getenv("PROJECT_ROOT"), "docker", "docker-compose.yml")
	// Call the function to spin up containers
	DockerManager = docker.GetInstance()
	err := DockerManager.InitializeStack(ctx, composeFile)
	if err != nil {
		log.Fatalf("Error initializing stack: %v", err)
	}
	stack := DockerManager.GetStack()

	obaContainer, err = DockerManager.GetOBAServerContainer(ctx)
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

	// init the OBA client
	log.Printf("Initializing OBA client with base URL: %s\n", "http://localhost:8085/onebusaway-api-webapp")
	obaClient = onebusaway.NewClient(
		option.WithAPIKey("TEST"),
		option.WithBaseURL("http://localhost:8085/onebusaway-api-webapp"),
	)

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
