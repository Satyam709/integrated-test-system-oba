package static_tests

import (
	"context"
	"testing"
	"time"

	onebusaway "github.com/OneBusAway/go-sdk"
	timec "github.com/Satyam709/integrated-test-system-oba/internal/time"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ArrivalsAndDeparturesTest struct {
	testTitle   string
	stopId      string
	params      onebusaway.ArrivalAndDepartureListParams
	currentTime int64

	expectedError                         string
	expectedNumberOfArrivalsAndDepartures int
	expectedArrivalsAndDepartures         []ExpectedArrivalsAndDeparture
}

type ExpectedArrivalsAndDeparture struct {
	tripId               string
	numberOfStopsAway    int64
	routeId              string
	scheduledArrivalTime int64
	distanceFromStop     float64
	stopSequence         int64
}

// add test data
func addTestData() []ArrivalsAndDeparturesTest {

	testCase1 := ArrivalsAndDeparturesTest{
		testTitle: "Many Arrivals and Departures ",
		params: onebusaway.ArrivalAndDepartureListParams{
			MinutesAfter:  onebusaway.Int(64),
			MinutesBefore: onebusaway.Int(5),
		},
		stopId:                                "1_600",
		currentTime:                           1747666800000,
		expectedError:                         "",
		expectedNumberOfArrivalsAndDepartures: 78,
		expectedArrivalsAndDepartures: []ExpectedArrivalsAndDeparture{
			{
				tripId:               "1_585676068",
				numberOfStopsAway:    -4,
				distanceFromStop:     -2369.6966623130174,
				routeId:              "1_102615",
				scheduledArrivalTime: 1747666502000,
				stopSequence:         4,
			},
			{
				tripId:               "1_628189348",
				numberOfStopsAway:    4,
				distanceFromStop:     1553.8917440443547,
				routeId:              "1_102615",
				scheduledArrivalTime: 1747667102000,
				stopSequence:         4,
			},
			{
				tripId:               "1_693954658",
				numberOfStopsAway:    20,
				distanceFromStop:     9635.657916198878,
				routeId:              "1_100169",
				scheduledArrivalTime: 1747668669000,
				stopSequence:         6,
			},
		},
	}

	return []ArrivalsAndDeparturesTest{
		testCase1,
	}
}

func TestArrivalsAndDepartures(t *testing.T) {
	tests := addTestData()

	for _, test := range tests {
		t.Run(test.testTitle, func(t *testing.T) {
			validate(test, t)
		})
	}
}

func validate(test ArrivalsAndDeparturesTest, t *testing.T) {

	// set the time
	err := timec.SetFakeTime(test.currentTime)
	assert.NoError(t, err, "Error setting fake time for arrivals and departures test: %s", test.testTitle)

	context, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()

	res, err := obaClient.ArrivalAndDeparture.List(context, test.stopId, test.params)

	if test.expectedError != "" {
		assert.NoError(t, err, "Error in arrivals departures api response for test: %s", test.testTitle)
		assert.EqualError(t, err, test.expectedError, "Expected error for test: %s", test.testTitle)
		return
	}

	require.NotEmpty(t, res, "Response should not be empty for test: %s", test.testTitle)

	assert.Equal(t, res.Code, int64(200), "Expected status code 200 for test: %s", test.testTitle)
	gotArrivals := res.Data.Entry.ArrivalsAndDepartures
	assert.Equal(t, test.expectedNumberOfArrivalsAndDepartures, len(gotArrivals), "Expected number of arrivals and departures for test: %s", test.testTitle)

	for _, expected := range test.expectedArrivalsAndDepartures {
		found := false
		for _, got := range gotArrivals {
			if expected.tripId == got.TripID &&
				expected.routeId == got.RouteID &&
				expected.numberOfStopsAway == got.NumberOfStopsAway &&
				expected.distanceFromStop == got.DistanceFromStop &&
				expected.scheduledArrivalTime == got.ScheduledArrivalTime &&
				expected.stopSequence == got.StopSequence {
				found = true
				break
			}
		}
		assert.True(t, found, "Expected trip ID %s not found in arrivals and departures for test: %s", expected.tripId, test.testTitle)
	}
}
