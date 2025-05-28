package docker

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type OBAContainer struct {
	Container *testcontainers.DockerContainer
}

func (o *OBAContainer) BuildBundle(ctx context.Context) error {
	log.Printf("Building GTFS bundle...")
	if o.Container == nil {
		return fmt.Errorf("OBA container is not initialized")
	}
	cmd := []string{"sh", "-c", "cd /bundle && ./build_bundle.sh"}
	_, _, err := o.Container.Exec(ctx, cmd)
	if err != nil {
		return fmt.Errorf("failed to execute build_bundle.sh: %w", err)
	}
	log.Printf("build_bundle.sh executed successfully")
	return nil
}

func (o *OBAContainer) StartServer(ctx context.Context) error {
	log.Printf("Starting OBA server...")
	if o.Container == nil {
		return fmt.Errorf("OBA container is not initialized")
	}
	cmd := []string{"sh", "-c", "/usr/local/tomcat/bin/catalina.sh start"}
	_, _, err := o.Container.Exec(ctx, cmd)
	if err != nil {
		return fmt.Errorf("failed to start OBA server: %w", err)
	}
	log.Printf("OBA server start command executed successfully")
	return nil
}

func (o *OBAContainer) WaitForServerReady(ctx context.Context) error {
	log.Printf("Waiting for OBA server to start...")
	if o.Container == nil {
		return fmt.Errorf("OBA container is not initialized")
	}

	err := wait.
		ForHTTP("/onebusaway-api-webapp/api/where/config.json?key=TEST").
		WithPort("8080/tcp").
		WithStartupTimeout(180*time.Second).
		WaitUntilReady(ctx, o.Container)

	if err != nil {
		return fmt.Errorf("OBA server did not start in time: %w", err)
	}
	log.Printf("OBA server is ready and running on port 8080")
	return nil
}

func (o *OBAContainer) StopServer(ctx context.Context) error {
	log.Printf("Stopping OBA server...")
	if o.Container == nil {
		return fmt.Errorf("OBA container is not initialized")
	}
	cmd := []string{"sh", "-c", "/usr/local/tomcat/bin/catalina.sh stop"}
	_, _, err := o.Container.Exec(ctx, cmd)
	if err != nil {
		return fmt.Errorf("failed to stop OBA server: %w", err)
	}
	log.Printf("OBA server stopped successfully")
	return nil
}

func (o *OBAContainer) Restart(ctx context.Context) error {
	log.Printf("Restarting OBA server...")
	if o.Container == nil {
		return fmt.Errorf("OBA container is not initialized")
	}
	err := o.StopServer(ctx)
	if err != nil {
		return fmt.Errorf("failed to stop OBA server: %w", err)
	}
	err = o.StartServer(ctx)
	if err != nil {
		return fmt.Errorf("failed to start OBA server: %w", err)
	}
	log.Printf("OBA server restarted successfully")
	return nil
}