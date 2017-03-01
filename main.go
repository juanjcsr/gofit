package main

import "github.com/juanjcsr/gofit/fitbit"

func main() {
	// client, _ := fitbit.NewFitbitClient()
	// // resp, err := client.Activities.GetActivitySummary("", "2017-02-13")
	// resp, err := client.Heart.GetIntradayHeartData("-", "today", "1sec", "00:00", "00:30")
	// if err != nil {
	// 	fmt.Printf("error getting user: %s", err)
	// }
	// fmt.Printf("user: %v", resp)
	fitbit.NewFitbitServer(8080)
}
