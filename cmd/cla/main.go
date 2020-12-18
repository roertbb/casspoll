package main

import (
	"fmt"
	"log"

	"github.com/manifoldco/promptui"
	"github.com/roertbb/casspoll/poll"
	cass "github.com/roertbb/casspoll/repo/cassandra"
)

const listPolls string = "List polls"
const createNewPoll string = "Create new poll"
const votePoll string = "Vote"
const getPollResults string = "See results"

func menu(s *poll.PollService) {
	prompt := promptui.Select{
		Label: "Select option",
		Items: []string{
			listPolls,
			createNewPoll,
		},
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	clearScreen()

	switch result {
	case listPolls:
		list(s)
	case createNewPoll:
		fmt.Println("TODO: create new poll")
	}
}

func list(s *poll.PollService) {
	polls, err := (*s).GetActivePolls()
	if err != nil {
		log.Fatal(err)
		return
	}

	pollList := []string{}
	pollMap := map[string]poll.Poll{}
	for _, poll := range *polls {
		id := fmt.Sprintf("[%s] %s", poll.ID.String(), poll.Description)
		pollList = append(pollList, id)
		pollMap[id] = poll
	}

	prompt := promptui.Select{
		Label: "Select poll",
		Items: pollList,
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	clearScreen()

	selectedPoll := pollMap[result]
	pollDetails(s, &selectedPoll)
}

func pollDetails(s *poll.PollService, p *poll.Poll) {
	answers, err := (*s).GetAnswers(p.ID)
	if err != nil {
		log.Fatal(err)
		return
	}

	pollTypeToString := map[poll.PollType]string{
		0: "Single-Choice",
		1: "Multiple-Choice",
	}

	fmt.Println(fmt.Sprintf("Title: %s", p.Title))
	fmt.Println(fmt.Sprintf("Description: %s", p.Description))
	fmt.Println(fmt.Sprintf("DueTime: %s", p.DueTime))
	fmt.Println(fmt.Sprintf("PollType: %s", pollTypeToString[p.PollType]))
	fmt.Println("--------")
	for _, answer := range *answers {
		fmt.Println(answer.Text)
	}

	prompt := promptui.Select{
		Label: "Select option",
		Items: []string{votePoll, getPollResults},
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	clearScreen()

	switch result {
	case votePoll:
		fmt.Println("TODO: vote")
	case getPollResults:
		fmt.Println("TODO: get poll results")
	}
}

func clearScreen() {
	print("\033[H\033[2J")
}

func main() {
	repo, _ := cass.NewCassandraRepo()
	service := poll.NewPollService(repo)

	menu(&service)
}
