package static_tests

import (
	"context"
	"log"
	"sort"
	"testing"
	"time"

	onebusaway "github.com/OneBusAway/go-sdk"
	"github.com/Satyam709/integrated-test-system-oba/internal/timecontroller"
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

	testCase2 := ArrivalsAndDeparturesTest{
		testTitle: "No Arrivals and Departures (very short time range)",
		params: onebusaway.ArrivalAndDepartureListParams{
			MinutesAfter:  onebusaway.Int(1),
			MinutesBefore: onebusaway.Int(2),
		},
		stopId:                                "1_26080",
		currentTime:                           1747927740000,
		expectedError:                         "",
		expectedNumberOfArrivalsAndDepartures: 0,
		expectedArrivalsAndDepartures:         []ExpectedArrivalsAndDeparture{},
	}

	// trip "1_729861698" is scheduled to arrive at stop "1_26080" at 1747869094000 on a service day
	testCase3 := ArrivalsAndDeparturesTest{
		testTitle: "Very few Arrivals and Departures",
		params: onebusaway.ArrivalAndDepartureListParams{
			MinutesBefore: onebusaway.Int(5),
			MinutesAfter:  onebusaway.Int(35),
		},
		stopId:                                "1_26080",
		currentTime:                           1747695600000,
		expectedError:                         "",
		expectedNumberOfArrivalsAndDepartures: 3,
		expectedArrivalsAndDepartures: []ExpectedArrivalsAndDeparture{
			{
				tripId:               "1_729861698",
				distanceFromStop:     -353.7095103595527,
				numberOfStopsAway:    -1,
				routeId:              "1_100254",
				scheduledArrivalTime: 1747695515000,
				stopSequence:         5,
			},
		},
	}

	// trip "1_729861698" should not be present on a non service day i.e. wednesday
	testCase4 := ArrivalsAndDeparturesTest{
		testTitle: "No Arrivals of trip 1_729861698 on non service day",
		params: onebusaway.ArrivalAndDepartureListParams{
			MinutesBefore: onebusaway.Int(5),
			MinutesAfter:  onebusaway.Int(35),
		},
		stopId:                                "1_26080",
		currentTime:                           1747868400000,
		expectedError:                         "",
		expectedNumberOfArrivalsAndDepartures: 2,
		expectedArrivalsAndDepartures: []ExpectedArrivalsAndDeparture{
			{
				distanceFromStop:     3718.4339441993798,
				numberOfStopsAway:    11,
				routeId:              "1_100254",
				tripId:               "1_724948468",
				scheduledArrivalTime: 1747869094000,
				stopSequence:         30,
			},
			{
				distanceFromStop:     8270.515801228234,
				numberOfStopsAway:    24,
				routeId:              "1_100254",
				tripId:               "1_724948328",
				scheduledArrivalTime: 1747870054000,
				stopSequence:         30,
			},
		},
	}

	testCases := []ArrivalsAndDeparturesTest{
		testCase1, testCase2, testCase3, testCase4,
	}

	// Sort test cases by currentTime to ensure consistent ordering
	// fixes server frizziness(occurs if we go back in past)
	sort.Slice(testCases, func(i, j int) bool {
		return testCases[i].currentTime < testCases[j].currentTime
	})

	return testCases
}

func TestArrivalsAndDepartures(t *testing.T) {
	tests := addTestData()
	restartObaServer()

	for _, test := range tests {
		t.Run(test.testTitle, func(t *testing.T) {
			validate(test, t)
		})
	}
}

func validate(test ArrivalsAndDeparturesTest, t *testing.T) {

	// set the time
	err := timecontroller.SetFakeTime(test.currentTime)
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
		if !found {
			log.Println(res.JSON.RawJSON())
		}
		assert.True(t, found, "Expected trip ID %s not found in arrivals and departures for test: %s", expected.tripId, test.testTitle)
	}
}
