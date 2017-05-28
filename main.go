package main

import (
	"log"
)

func main() {
	// client, _ := fitbit.NewFitbitClient()
	// // resp, err := client.Activities.GetActivitySummary("", "2017-02-13")
	// resp, err := client.Heart.GetIntradayHeartData("-", "today", "1sec", "00:00", "00:30")
	// if err != nil {
	// 	fmt.Printf("error getting user: %s", err)
	// }
	db, err := DB()
	if err != nil {
		log.Fatalf("could not open DB: %v", err)
	}
	db.AutoMigrate(&User{}, &Token{})
	NewFitbitServer(3000)
}
