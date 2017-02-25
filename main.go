package main

import (
	"fmt"

	"github.com/juanjcsr/gofit/fitbit"
)

func main() {
	// fitbit.NewFitbitClient()
	client, _ := fitbit.NewFitbitClient()
	resp, err := client.Activities.GetActivitySummary("", "2017-02-13")
	if err != nil {
		fmt.Printf("error getting user: %s", err)
	}
	fmt.Printf("user: %v", resp)
}
