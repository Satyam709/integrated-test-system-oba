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
	"github.com/Satyam709/integrated-test-system-oba/internal/time"
)

var (
	DockerManager *docker.DockerManager
	// ctx           context.Context
	obaClient    *onebusaway.Client
	obaContainer *docker.OBAContainer
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
	defer cleanUP()

	obaContainer, err = DockerManager.GetOBAServerContainer(ctx)
	if err != nil {
		log.Fatalf("Error getting OBA server container: %v", err)
	}

	// build_bundle.sh log
	err = obaContainer.BuildBundle(ctx)
	if err != nil {
		log.Fatalf("Error building GTFS bundle: %v", err)
	}

	// init the OBA client
	log.Printf("Initializing OBA client with base URL: %s\n", "http://localhost:8085/onebusaway-api-webapp")
	obaClient = onebusaway.NewClient(
		option.WithAPIKey("TEST"),
		option.WithBaseURL("http://localhost:8085/onebusaway-api-webapp"),
	)

	// set initial time
	err = time.SetFakeTime(1746082800000) // 2025-05-1 00:00:00
	if err != nil {
		log.Fatalf("Error setting initial fake time: %v", err)
	}

	// Run the tests
	fmt.Printf("Running static tests...\n")
	exitCode := m.Run()
	fmt.Printf("exitCode: %v\n", exitCode)
}

func cleanUP() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := DockerManager.Cleanup(ctx)
	if err != nil {
		log.Fatalf("Error cleaning up Docker stack: %v", err)
	}
	log.Printf("Docker stack cleaned up successfully")
}

func restartObaServer() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if obaContainer == nil {
		log.Fatalf("OBA container is not initialized")
	}
	err := obaContainer.Restart(ctx)
	if err != nil {
		log.Fatalf("Error restarting OBA server: %v", err)
	}
	// Wait for port 8080 (inside container)
	err = obaContainer.WaitForServerReady(ctx)
	if err != nil {
		log.Fatalf("Error waiting for OBA server to be ready: %v", err)
	}
}
