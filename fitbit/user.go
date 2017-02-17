package fitbit

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type UserHolder struct {
	User User `json:"user"`
}

// User represents a Fitbit user
type User struct {
	Age                    int
	AutoStrideEnabled      bool
	Avatar                 string
	Avatar150              string
	AverageDailySteps      int16
	ClockTimeDisplayFormat string
	Corporate              bool
	CorporateAdmin         bool
	Country                string
	DateOfBirth            string
	DisplayName            string
	DisplayNameSetting     string
	DistanceUnit           string
	ID                     string `json:"encodedID"`
	//Features
	FoodsLocale             string
	FullName                string
	Gender                  string
	GlucoseUnit             string
	Heigth                  int16
	HeightUnit              string
	Locale                  string
	MemberSince             string
	OffsetFromUTCMillis     int64
	StartDayOfWeek          string
	StrideLengthRunning     float32
	StrideLengthRunningType string
	StrideLengthWalking     float32
	StrideLengthWalkingType string
	SwimUnit                string
	Timezone                string
	WaterUnit               string
	WaterUnitName           string
	Weight                  int
	WeightUnit              string
}

const userEndpoint = "https://api.fitbit.com/1/user/-/profile.json"

// UserService holds the requests for the user endpoints
type UserService struct {
	client *http.Client
}

func newUserService(client *http.Client) *UserService {
	var userService = new(UserService)
	userService.client = client
	return userService
}

func (u *UserService) GetCurrentUser() (*User, error) {
	resp, err := u.client.Get(userEndpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// json.NewDecoder(resp.Body)
	var uh UserHolder
	// var user User
	if err := json.NewDecoder(resp.Body).Decode(&uh); err != nil {
		return nil, fmt.Errorf("user endpoint: %v", err)
	}
	// respD, err := ioutil.ReadAll(resp.Body)
	// fmt.Printf("err: %v", Re)
	fmt.Printf("user: %v", uh)

	// return &uh.User, nil
	return &uh.User, nil
}
