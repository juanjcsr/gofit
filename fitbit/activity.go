package fitbit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type ActivitiesHolder struct {
	Activities []Activity `json:"activities"`
	Goals      Goals      `json:"goals"`
	Summary    Summary    `json:"summary"`
}

type Source struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	URL  string `json:"url"`
}

type ActivitiesLogHolder struct {
	Activities []ActivityLog `json:"activities"`
}

type ActivityLog struct {
	ActiveDuration       int32           `json:"activeDuration"`
	ActivityLevels       []ActivityLevel `json:"activityLevel"`
	ActivityName         string          `json:"activityName"`
	ActivityTypeID       int             `json:"activityTypeId"`
	AverageHeartRate     int             `json:"averageHeartRate"`
	Calories             int             `json:"calories"`
	CaloriesLink         string          `json:"caloriesLink"`
	CustomHeartRateZones []HeartRateZone `json:"customHeartRateZones"`
	Distance             float32         `json:"distance"`
	DistanceUnit         string          `json:"distanceUnit"`
	HeartRateLink        string          `json:"heartRateLink"`
	HeartRateZones       []HeartRateZone `json:"heartRateZones"`
	LastModified         string          `json:"lastModified"`
	LogID                int64           `json:"logId"`
	LogType              string          `json:"logType"`
	OriginalDuration     int             `json:"originalDuration"`
	OriginalStartTime    time.Time       `json:"originalStartTime"`
	Pace                 float64         `json:"pace"`
	Source               Source          `json:"source"`
	Speed                float64         `json:"speed"`
	StartTime            time.Time       `json:"startTime"`
	TcxLink              string          `json:"tcxLink"`
}

type ActivityLevel struct {
	Minutes int    `json:"minutes"`
	Name    string `json:"name"`
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
	activitiesLogEndpoint  = "https://api.fitbit.com/1/user/%s/activities/list.json?sort=%s&limit=%d&offset=0%s"
)

func newActivityService(authedClient *http.Client) *ActivityService {
	var theActivityService = new(ActivityService)
	theActivityService.authedClient = authedClient
	return theActivityService
}

func (h *ActivityService) GetTCXData(activityLog ActivityLog) (string, error) {
	resp, err := h.authedClient.Get(activityLog.TcxLink)
	if err != nil {
		return "", fmt.Errorf("could not get tcx data: %v", err)
	}
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	return string(bodyBytes), nil
}

func (a *ActivityService) GetActivityLog(user string, beforeDate string, afterDate string, limit int) (*ActivitiesLogHolder, error) {
	var theDate string
	var sort string
	if user == "" {
		user = "-"
	}
	if beforeDate != "" && afterDate == "" {
		theDate = "&beforeDate=" + beforeDate
		sort = "desc"
	} else {
		theDate = "&afterDate=" + afterDate
		sort = "asc"
	}

	url := fmt.Sprintf(activitiesLogEndpoint, user, sort, limit, theDate)
	resp, err := a.authedClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("could not get activities log response for %s: %v", url, err)
	}
	defer resp.Body.Close()
	var ah ActivitiesLogHolder
	if err := json.NewDecoder(resp.Body).Decode(&ah); err != nil {
		return nil, fmt.Errorf("could not decode activities: %v", err)
	}
	return &ah, nil
}

func (a *ActivityService) GetActivitySummary(user string, date string) (*ActivitiesHolder, error) {
	var url string
	if user != "" {
		url = fmt.Sprintf(activitiesBaseEndpoint, user, date)
	} else {
		url = fmt.Sprintf(activitiesBaseEndpoint, "-", date)
	}
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

func (a *ActivitiesHolder) ToJson() []byte {
	jsonExp, _ := json.Marshal(a)
	return jsonExp
}

func (a *ActivitiesLogHolder) ToJson() []byte {
	jsonExp, _ := json.Marshal(a)
	return jsonExp
}
