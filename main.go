package main

import "github.com/juanjcsr/gofit/fitbit"
import "fmt"

func main() {
	client, _ := fitbit.NewFitbitClient()
	resp, err := client.User.GetUser("")
	if err != nil {
		fmt.Printf("error getting user: %s", err)
	}
	fmt.Printf("user: %v", resp)
}
