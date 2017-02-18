package fitbit

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ActivitiesHolder struct {
	Activities []Activity `json:"activities"`
	Goals      Goals      `json:"goals"`
	Summary    Summary    `json:"summary"`
}

type Activity struct {
	ID                 int32  `json:"activityId"`
	ActivityParentID   int32  `json:"activityId"`
	ActivityParentName string `json:"activityParentName"`
	Calories           int32  `json:"activityId"`
	Description        string `json:"description"`
	Distance           int32  `json:"distance"`
	Duration           int32  `json:"duration"`
	HasStartTime       bool   `json:"hasStartTime"`
	IsFavorite         bool   `json:"isFavorite"`
	LastModified       string `json:"lastModified"`
	LogID              int64  `json:"logId"`
	Name               string `json:"name"`
	StartDate          string `json:"startDate"`
	StartTime          string `json:"startTime"`
	Steps              int32  `json:"steps"`
}

type Goals struct {
	ActiveMinutes int16   `json:"activeMinutes"`
	CaloriesOut   int16   `json:"caloriesOut"`
	Distance      float32 `json:"distance"`
	Floors        int16   `json:"floors"`
	Steps         int16   `json:"steps"`
}

type Summary struct {
	ActiveScore         int16      `json:"activeScore"`
	ActivityCalories    int16      `json:"activityCalories"`
	CaloriesBMR         int16      `json:"caloriesBMR"`
	CaloriesOut         int16      `json:"caloriesOut"`
	Distances           []Distance `json:"distances"`
	Elevation           float32    `json:"elevation"`
	FairlyActiveMinutes int16      `json:"fairlyActiveMinutes"`
	Floors              int16      `json:"floors"`
	// HEART RATE ZONES
	LightlyActiveMinutes int16 `json:"lightlyActiveMinutes"`
	MarginalCalories     int16 `json:"marginalCalories"`
	RestingHeartRate     int16 `json:"restingHeartRate"`
	SedentaryMinutes     int16 `json:"sedentaryMinutes"`
	Steps                int32 `json:"steps"`
	VeryActiveMinutes    int16 `json:"veryActiveMinutes"`
}

type Distance struct {
	Activity string  `json:"activity"`
	Distance float32 `json:"distance"`
}

type ActivityService struct {
	authedClient *http.Client
}

const (
	activitiesBaseEndpoint = "https://api.fitbit.com/1/user/%s/activities/date/%s.json"
)

func newActivityService(authedClient *http.Client) *ActivityService {
	var theActivityService = new(ActivityService)
	theActivityService.authedClient = authedClient
	return theActivityService
}

func (a *ActivityService) GetActivitySummary(user string) (*ActivitiesHolder, error) {
	var url string
	if user != "" {
		url = fmt.Sprintf(activitiesBaseEndpoint, user, "2017-02-16")
	} else {
		url = fmt.Sprintf(activitiesBaseEndpoint, "-", "2017-02-13")
	}
	fmt.Printf("URL: %s\n", url)
	resp, err := a.authedClient.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("user service: could not get user activities: %s", err)
	}
	var ah ActivitiesHolder
	if err := json.NewDecoder(resp.Body).Decode(&ah); err != nil {
		return nil, fmt.Errorf("user service: could not parse user activities: %s", err)
	}
	return &ah, nil
}
