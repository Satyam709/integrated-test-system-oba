package static_tests

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/Satyam709/integrated-test-system-oba/docker_utils"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestMain(m *testing.M) {
	// Set up the environment for static tests
	// This could include setting up environment variables, initializing configurations, etc.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Call the function to spin up containers
	stack := docker_utils.SpinUpContainers(ctx, "../../docker/docker-compose.yml")

	// get the container and build bundle and start server
	oba_container, err := stack.ServiceContainer(ctx, "oba_server")
	if err != nil {
		fmt.Printf("Error getting OBA server container: %v\n", err)
		return
	}

	// build_bundle.sh log
	_, _, err = oba_container.Exec(ctx, []string{"sh", "-c", "cd /bundle && ./build_bundle.sh"})
	if err != nil {
		log.Fatalf("Error executing build_bundle.sh: %v", err)
	}
	fmt.Printf("build_bundle.sh executed successfully\n")
	_, _, err = oba_container.Exec(ctx, []string{"sh", "-c", "/usr/local/tomcat/bin/catalina.sh start"})
	if err != nil {
		log.Fatalf("Error starting OBA server: %v", err)
	}
	fmt.Printf("Server start cmd executed\n")

	// Wait for port 8080 (inside container)
	log.Printf("waiting for server to start ...")
	err = wait.
		ForHTTP("/onebusaway-api-webapp/api/where/config.json?key=TEST").
		WithPort("8080/tcp").
		WithStartupTimeout(180*time.Second).
		WaitUntilReady(ctx, oba_container)

	if err != nil {
		log.Fatalf("OBA server not ready: %v", err)
	} else {
		fmt.Printf("OBA server is ready and running on port 8080\n")
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
