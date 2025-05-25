package static_tests

import (
	"encoding/json"
	"log"
	"net/http"
	"testing"
	
	"github.com/stretchr/testify/assert"
)

func TestISWorking(t *testing.T) {
	assert.True(t, true, "True is true!")
}


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
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	log.Printf("Response: %v\n", response)
	assert.NotNil(t, response, "Response should not be nil")
	assert.Contains(t, response, "currentTime", "Response should contain 'currentTime' key")
	log.Printf("Current time response: %v\n", response["currentTime"])
}
