package static_tests

import (
	"context"
	"log"
	"testing"
	"time"

	timec "github.com/Satyam709/integrated-test-system-oba/internal/time"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCurrentTime(t *testing.T) {
	log.Println("Starting test for current time endpoint...")
	
	expectedTime :=int64(1747783080000)
	// Mock the current time endpoint
	err := timec.SetFakeTime(expectedTime)
	assert.NoError(t, err, "Expected no error when setting fake time")

	// restart the oba server
	restartObaServer()

	contextWithTimeout, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
	defer cancel()

	client := obaClient
	currentTime, err := client.CurrentTime.Get(contextWithTimeout)
	assert.NoError(t, err, "Expected no error when getting current time")
	require.NotNil(t, currentTime, "Expected current time response to be non-nil")
	assert.Equal(t,currentTime.Code, int64(200), "Expected status code 200")
	// fmt.Printf("Response : %+v\n", currentTime)
	assert.Equal(t, expectedTime, currentTime.CurrentTime, "Expected current time to match the mocked time")
}