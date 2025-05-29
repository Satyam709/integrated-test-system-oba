# Integrated Test System for OneBusAway API

A test harness for validating OneBusAway API server behavior using known GTFS data inputs.

## Prerequisites

- Docker
- Docker Compose
- Go 1.24+

## Setup

1. Clone the repository:
```bash
git clone https://github.com/Satyam709/integrated-test-system-oba.git
```

2. Download OBA artifacts:
```bash
cd oba-artifacts
./retrieve-oba-artifacts.sh
```

## Running Tests

1. Build the test environment image:
```bash
docker build -t oba_server_testing_image_v1 -f docker/Dockerfile .
```

2. Run the test suite:
```bash
go test ./tests/static_tests/...
```

## Test Environment

- OBA Server: http://localhost:8085
- MySQL Database: localhost:3310
  - Database: oba_database
  - User: oba_user
  - Password: oba_password

## Project Structure

- `/docker`: Docker configuration and scripts
- `/testdata`: GTFS test data
- `/tests`: Test suites
- `/internal`: Internal packages for test infrastructure