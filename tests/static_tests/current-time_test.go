package static_tests

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	timec "github.com/Satyam709/integrated-test-system-oba/internal/time"
)

func TestCurrentTime(t *testing.T) {
	log.Println("Starting test for current time endpoint...")
	
	expectedTime :=int64(1747539000000)
	// Mock the current time endpoint
	err := timec.SetFakeTime(expectedTime)
	assert.NoError(t, err, "Expected no error when setting fake time")

	client := obaClient
	currentTime, err := client.CurrentTime.Get(context.TODO())
	assert.Equal(t,currentTime.Code, int64(200), "Expected status code 200")
	// fmt.Printf("Response : %+v\n", currentTime)
	assert.NoError(t, err, "Expected no error when getting current time")
	assert.Equal(t, expectedTime, currentTime.CurrentTime, "Expected current time to match the mocked time")
}