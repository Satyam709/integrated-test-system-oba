package static_tests

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"testing"
	"time"

	onebusaway "github.com/OneBusAway/go-sdk"
	"github.com/OneBusAway/go-sdk/option"
	"github.com/stretchr/testify/assert"
	timec "github.com/Satyam709/integrated-test-system-oba/internal/time"
)

func TestCurrentTime(t *testing.T) {
	// make a req to http://localhost:8085/onebusaway-api-webapp/api/where/current-time.json?key=TEST
	log.Println("Testing current time endpoint...")
	resp, err := http.Get("http://localhost:8085/onebusaway-api-webapp/api/where/current-time.json?key=TEST")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()
	log.Printf("Response status code: %d\n", resp.StatusCode)
	assert.Equal(t, 200, resp.StatusCode, "Expected status code 200, got %d", resp.StatusCode)
	var response map[string]any
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	log.Printf("Response: %v\n", response)
	assert.NotNil(t, response, "Response should not be nil")
	assert.Contains(t, response, "currentTime", "Response should contain 'currentTime' key")
	log.Printf("Current time response: %v\n", response["currentTime"])
}

func TestCurrentTimeWithSDK(t *testing.T) {
	expectedTime :=int64(1747539000000)
	// Mock the current time endpoint
	err := timec.SetFakeTime(expectedTime)
	
	assert.NoError(t, err, "Expected no error when setting fake time")

	log.Printf("waiting for fake time to be set: %d\n", expectedTime)
	time.Sleep(10 * time.Second) // Ensure the fake time is set before making the request
	log.Println("Testing current time endpoint with SDK...")
	client := onebusaway.NewClient(
		option.WithAPIKey("TEST"),
		option.WithBaseURL("http://localhost:8085/onebusaway-api-webapp"),
	)

	currentTime, err := client.CurrentTime.Get(context.TODO())
	// assert.Equal(t,currentTime.Code, onebusaway.Int(200), "Expected status code 200")
	// fmt.Printf("Response : %+v\n", currentTime)
	assert.NoError(t, err, "Expected no error when getting current time")
	assert.Equal(t, expectedTime, currentTime.CurrentTime, "Expected current time to match the mocked time")
}