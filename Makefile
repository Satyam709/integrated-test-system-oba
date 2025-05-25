.PHONY: test-all test-static test-realtime clean

# env variables
export PROJECT_ROOT=$(shell pwd)

# Default target
all: test-all

# Run all tests in sequence
test-all:
	make test-static && make test-realtime

# Run only static tests
test-static:
	go test -v ./tests/static_tests/...

# Run only realtime tests
test-realtime:
	go test ./tests/realtime_tests/...

test-time_controller:
	go test ./internal/time/...
# Clean up any temporary files or artifacts
clean:
	rm -f tests/static_tests/*.test
	rm -f tests/realtime_tests/*.test