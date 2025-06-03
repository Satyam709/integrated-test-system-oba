.PHONY: test-all test-static test-realtime clean

# env variables
export PROJECT_ROOT=$(shell pwd)
export TIME_CONFIG_PATH=$(shell pwd)/internal/timecontroller/faketime.cfg

# Default target
all: test-all

# Run all tests in sequence
test-all:
	make test-static && make test-realtime

# Run only static tests
test-static:
	go test -count=1 -v ./tests/static_tests/...

# Run only realtime tests
test-realtime:
	go test -count=1 -v ./tests/realtime_tests/...

test-time_controller:
	go test -count=1 -v ./internal/timecontroller/...
# Clean up any temporary files or artifacts
clean:
	rm -f tests/static_tests/*.test
	rm -f tests/realtime_tests/*.test