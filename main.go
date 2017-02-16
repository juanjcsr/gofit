package main

import "github.com/juanjcsr/gofit/fitbit"
import "fmt"
import "io/ioutil"

func main() {
	client, _ := fitbit.NewFitbitClient()
	resp, err := client.Client.Get("https://api.fitbit.com/1/user/-/profile.json")
	if err != nil {
		fmt.Printf("error getting user: %s", err)
	}
	defer resp.Body.Close()
	resultByte, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("user: %v", string(resultByte[:]))
}
