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
	TopBadges               []Badge
	WaterUnit               string
	WaterUnitName           string
	Weight                  int
	WeightUnit              string
}

// Badge represents a Badge from the User
type Badge struct {
	BadgeGradientEndColor   string `json:"badgeGradientEndColor"`
	BadgeGradientStartColor string `json:"badgeGradientStartColor"`
	BadgeType               string `json:"badgeType"`
	Category                string `json:"category"`
	Cheers                  []interface{}
	DateTime                string `json:"dateTime"`
	Description             string `json:"description"`
	BadgeID                 string `json:"encodedId"`
	Image100px              string `json:"image100px"`
	Image125px              string `json:"image100px"`
	Image300px              string `json:"image100px"`
	Image50px               string `json:"image100px"`
	Image75x                string `json:"image100px"`
	MarketingDescription    string `json:"marketingDescription"`
	MobileDescription       string `json:"mobileDescription"`
	Name                    string `json:"name"`
	ShareImage              string `json:"shareImage"`
	ShareText               string `json:"shareText"`
	ShortDescription        string `json:"shortDescription"`
	ShortName               string `json:"shortName"`
	TimesArchived           int16  `json:"timesArchived"`
	Value                   int16  `json:"value"`
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
