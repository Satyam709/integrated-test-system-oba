package timecontroller

import (
	"testing"
	"os"
)

func TestSetFakeTime(t *testing.T) {
	// Test timestamp: 2024-03-20 15:30:00 PST
	// (1710974400000 milliseconds since epoch)
	testTime := int64(1747539000000)

	err := SetFakeTime(testTime)
	if err != nil {
		t.Fatalf("SetFakeTime failed: %v", err)
	}

	// Get the current working directory
	configPath := os.Getenv("TIME_CONFIG_PATH")

	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read faketime.cfg: %v", err)
	}

	// Parse the written time
	writtenTime := string(content)
	expectedTime := "2025-05-17 20:30:00"

	// Compare the written time with expected
	if writtenTime != expectedTime {
		t.Errorf("Expected time %s, but got %s", expectedTime, writtenTime)
	}

	// Verify file permissions
	fileInfo, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("Failed to get file info: %v", err)
	}

	// Check if permissions are 0644 (rw-r--r--)
	expectedPerm := os.FileMode(0644)
	if fileInfo.Mode().Perm() != expectedPerm {
		t.Errorf("Expected file permissions %v, but got %v", expectedPerm, fileInfo.Mode().Perm())
	}
}