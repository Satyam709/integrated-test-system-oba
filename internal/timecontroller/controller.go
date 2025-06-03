package timecontroller

import (
	"fmt"
	"log"
	"os"
	"time"
)

// SetFakeTime converts milliseconds timestamp to Los Angeles timezone and writes to faketime.cfg
// Returns error if any operation fails
func SetFakeTime(milliseconds int64) error {
	// Convert milliseconds to time
	t := time.Unix(0, milliseconds*int64(time.Millisecond))

	// Load Los Angeles location
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		return fmt.Errorf("error loading timezone: %v", err)
	}

	// Convert time to Los Angeles timezone
	laTime := t.In(loc)

	formattedTime := laTime.Format("2006-01-02 15:04:05")
	log.Printf("Setting faketime to: %s", formattedTime)

	// Get the current working directory
	configPath := os.Getenv("TIME_CONFIG_PATH")

	err = os.WriteFile(configPath, []byte(formattedTime), 0644)
	if err != nil {
		return fmt.Errorf("error writing to config file: %v", err)
	}

	// wait 10 seconds to ensure libfaketime reads the new time
	time.Sleep(2 * time.Second)

	return nil
}
