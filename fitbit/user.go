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
	AboutMe                string `json:"aboutMe"`
	Avatar                 string
	Avatar150              string
	City                   string
	ClockTimeDisplayFormat string
	Country                string
	DateOfBirth            string
	DisplayName            string
	DistanceUnit           string
	EncodedID              string
	FoodsLocale            string
	FullName               string
	Gender                 string
	GlucoseUnit            string
	// Heigth                 string
	HeightUnit          string
	Locale              string
	MemberSince         string
	OffsetFromUTCMillis string
	StartDayOfWeek      string
	State               string
	StrideLengthRunning string
	StrideLengthWalking string
	Timezone            string
	WaterUnit           string
	// Weight                 string
	WeightUnit string
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
	uh := new(UserHolder)
	// user := new(User)
	if err := json.NewDecoder(resp.Body).Decode(uh); err != nil {
		return nil, fmt.Errorf("user endpoint: %v", err)
	}
	// fmt.Printf("USUARIO: %v", user)
	return &uh.User, nil
}
