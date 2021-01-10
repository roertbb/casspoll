package main

import (
	"fmt"

	"github.com/gocql/gocql"
)

func main() {
	// addresses := strings.Split(os.Getenv("ADDRESS"), ",")
	// if len(addresses) == 1 && addresses[0] == "" {
	// 	log.Fatal("ADDRESS env variable not specified")
	// 	os.Exit(1)
	// }

	address := "http://127.0.0.1:8080"

	createRandomPoll(address)

	polls, _ := getActivePolls(address)
	fmt.Println(polls)

	voterUUID, _ := gocql.RandomUUID()

	p := (*polls)[0]

	selectedAnswers, _ := voteRandomlyForPoll(address, &p, voterUUID)
	fmt.Println(selectedAnswers)

	res, _ := getPollResults(address, &p)
	for _, r := range *res {
		fmt.Println(r.Answer, r.VotesNo)
	}
}
