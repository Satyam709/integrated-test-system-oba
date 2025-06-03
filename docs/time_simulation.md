# Time Simulation Approach

The project implements time simulation functionality using libfaketime to control the system time within the OneBusAway test environment. Here's how it works:

## Components

### 1. Time Controller
Located in [Dockerfile](../internal/timecontroller/controller.go), the time controller provides:

- The file has function that:
  - Takes a Unix timestamp in milliseconds
  - Converts it to Los Angeles timezone
  - Writes the formatted time to faketime.cfg(on host)
  - Includes a 2-second delay to ensure libfaketime reads the new time

### 2. Docker Configuration

The [Time Controller](../docker/Dockerfile) sets up the libfaketime :
   - Clones from https://github.com/wolfcw/libfaketime
   - Builds and installs the library
   - Sets required environment variables:
     ```dockerfile
     ENV LD_PRELOAD=/usr/local/lib/faketime/libfaketime.so.1
     ENV FAKETIME_DONT_FAKE_MONOTONIC=1
     ```

## Usage

To simulate a specific time in your tests:

```go
err := timecontroller.SetFakeTime(milliseconds) // Pass Unix timestamp in milliseconds
if err != nil {
    // Handle error
}
```

The time change will affect all processes within the container that use system time calls, allowing for consistent time based testing.

Libfaketime caches the previous timestamp for the time specified in env Var `FAKETIME_CACHE_DURATION` currently this has been set to 1 sec to reduce the time for tests.

More info on libfaketime can be found at info on libfaketime can be found at [fakelibtime docs](https://github.com/wolfcw/libfaketime/blob/master/README)

## Best practises and limitations
- Although libfaketime sets the custom time for the processes within the container , it can fail for several reason
- Currently the biggest limitation of the time simulation approach is that the server running in a test container does not behave properly when we try to set a time past of what has been last most recent timestamp without restarting the server but we can go in future direction without restarting.

## Important Notes

1. Time changes are container-wide and affect all processes using system time calls
2. Times are always converted to Los Angeles timezone before being applied
3. A 2-second delay is included after time changes to ensure proper synchronization
4. The system uses a configuration file (faketime.cfg) that is bind to the /etc/faketimerc inside the test container to set time dynamically by libfaketime