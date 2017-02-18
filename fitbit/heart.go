package fitbit

type HeartRateZone struct {
	CaloriesOut float32 `json:"caloriesOut"`
	Max         int16   `json:"max"`
	Min         int16   `json:"min"`
	Minutes     int16   `json:"minutes"`
	Name        string  `json:"name"`
}

type HeartRateZones struct {
	CustomHeartRateZones []ActivitiesHeart `json:"customHeartRateZones"`
	HeartRateZone        []ActivitiesHeart `json:"heartRateZones"`
}

type HeartValues struct {
	DateTime string           `json:"dateTime"`
	Values   []HeartRateZones `json:"value"`
}

type ActivitiesHeart struct {
	HeartValues []HeartValues `json:"activities-Heart"`
}
