package fitbit

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HeartRateZone struct {
	CaloriesOut float32 `json:"caloriesOut"`
	Max         int16   `json:"max"`
	Min         int16   `json:"min"`
	Minutes     int16   `json:"minutes"`
	Name        string  `json:"name"`
}

type HeartRateZones struct {
	CustomHeartRateZones []HeartRateZone `json:"customHeartRateZones"`
	HeartRateZone        []HeartRateZone `json:"heartRateZones"`
	RestingHeartRate     int16           `json:"restingHeartRate"`
	Value                int16           `json:"value"`
	Datetime             string          `json:"dateTime"`
}
type HeartRateZonesIntraday struct {
	CustomHeartRateZones []HeartRateZone `json:"customHeartRateZones"`
	HeartRateZone        []HeartRateZone `json:"heartRateZones"`
	RestingHeartRate     int16           `json:"restingHeartRate"`
	Value                string          `json:"value"`
	Datetime             string          `json:"dateTime"`
}

type HRValue struct {
}

type HeartValues struct {
	DateTime string         `json:"dateTime"`
	Values   HeartRateZones `json:"value"`
}

type ActivitiesHeartHolder struct {
	HeartValues         []HeartValues           `json:"activities-heart"`
	HeartValuesIntraday ActivitiesHeartIntraday `json:"activities-heart-intraday"`
}

type ActivitiesHeartIntradayHolder struct {
	HeartValues         []HeartRateZonesIntraday `json:"activities-heart"`
	HeartValuesIntraday ActivitiesHeartIntraday  `json:"activities-heart-intraday"`
}

type HeartDataSet struct {
	Time  string `json:"time"`
	Value int16  `json:"value"`
}

type ActivitiesHeartIntraday struct {
	Dataset         []HeartDataSet `json:"dataset"`
	DatasetInterval int16          `json:"datasetInterval"`
	DatasetType     string         `json:"datasetType"`
}
type HeartService struct {
	client *http.Client
}

func newHeartService(authedClient *http.Client) *HeartService {
	var heartService = new(HeartService)
	heartService.client = authedClient
	return heartService
}

const (
	heartEndpoint          = "https://api.fitbit.com/1/user/%s/activities/heart/date/%s/%s.json"
	heartDateRangeEndpoint = "https://api.fitbit.com/1/user/%s/activities/heart/date/%s/%s.json"
	heartIntraEndpoint     = "https://api.fitbit.com/1/user/%s/activities/heart/date/%s/%s/%s.json"
	heartIntraDayEndpoint  = "https://api.fitbit.com/1/user/%s/activities/heart/date/%s/1d/%s/time/%s/%s.json"
)

func (h *HeartService) GetHeartData(userID string, date string, period string) (*ActivitiesHeartHolder, error) {
	var url string
	url = fmt.Sprintf(heartEndpoint, userID, date, period)
	resp, err := h.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var hh ActivitiesHeartHolder
	if err := json.NewDecoder(resp.Body).Decode(&hh); err != nil {
		return nil, fmt.Errorf("user endpoint: %v", err)
	}

	return &hh, nil
}

func (h *HeartService) GetIntradayHeartData(userID, date, detail, start, end string) (*ActivitiesHeartIntradayHolder, error) {
	// var url string
	url := fmt.Sprintf(heartIntraDayEndpoint, userID, date, detail, start, end)
	resp, err := h.client.Get(url)
	if err != nil {
		return nil, err
	}
	fmt.Println(url)
	defer resp.Body.Close()
	var intra ActivitiesHeartIntradayHolder
	if err := json.NewDecoder(resp.Body).Decode(&intra); err != nil {
		return nil, fmt.Errorf("heart endpoint: %v", err)
	}
	return &intra, nil
}
