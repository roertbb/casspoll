package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocql/gocql"
	"github.com/manifoldco/promptui"
	"github.com/roertbb/casspoll/poll"
	cass "github.com/roertbb/casspoll/repo/cassandra"
)

const (
	listPolls      = "List polls"
	createNewPoll  = "Create new poll"
	votePoll       = "Vote"
	getPollResults = "See results"
	singleChoice   = "Single-Choice"
	multipleChoice = "Multiple-Choice"
	quit           = "Quit"
	back           = "Back"
	sendVotes      = "Send votes"
)

var pollService poll.PollService
var voterID gocql.UUID
var pollTypeToString map[poll.PollType]string = map[poll.PollType]string{
	0: singleChoice,
	1: multipleChoice,
}

func menu() {
	_, result, _ := selectOption("Select option", []string{
		listPolls,
		createNewPoll,
		quit,
	})

	clearScreen()

	switch result {
	case listPolls:
		listActivePolls()
	case createNewPoll:
		createPoll()
	case quit:
		os.Exit(0)
	}
}

func createPoll() {
	var err error

	newPoll := poll.Poll{}

	if newPoll.Title, err = promptNonEmptyString("Enter title"); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	if newPoll.Description, err = promptNonEmptyString("Enter description"); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	if newPoll.DueTime, err = promptFutureDate("Enter due time (in YYYY-MM-DDTHH:MM:SS format)"); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	prompt := promptui.Select{
		Label: "Select type",
		Items: []string{singleChoice, multipleChoice},
	}

	choiceType, _, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}
	newPoll.PollType = poll.PollType(choiceType)

	fmt.Println("----------------")
	fmt.Println("Enter answers for the poll")
	fmt.Println("(type and confirm with <Enter> / leave empty and press <Enter> if all answers are added)")
	fmt.Println("----------------")

	answers := []string{}
	done := false
	for !done {
		answer, err := promptString("Enter answer")
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		if answer != "" {
			answers = append(answers, answer)
		} else {
			done = true
		}
	}

	reverse(answers)

	uuid, err := pollService.CreatePoll(&newPoll, &answers)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	clearScreen()

	fmt.Println("Successfully created poll!", uuid)
}

func listActivePolls() {
	polls, err := pollService.GetActivePolls()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	actions := []string{}
	for _, poll := range *polls {
		actions = append(actions, poll.Description)
	}
	actions = append(actions, back)

	idx, _, _ := selectOption("Select poll", actions)
	clearScreen()

	if idx < len(*polls) {
		pollDetails(&(*polls)[idx])
	}
}

func pollDetails(p *poll.Poll) {
	answers, err := pollService.GetAnswers(p.ID)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	logPollDetails(p)
	fmt.Println("Answers:")
	for _, answer := range *answers {
		fmt.Println(fmt.Sprintf("- %s", answer.Text))
	}

	_, result, _ := selectOption("Select option", []string{votePoll, getPollResults, back})

	clearScreen()

	switch result {
	case back:
		return
	case votePoll:
		vote(p, *answers)
	case getPollResults:
		getResults(p)
	}
}

func vote(p *poll.Poll, answers []poll.Answer) {
	votes := []gocql.UUID{}

	options := []string{}
	for _, answer := range answers {
		options = append(options, answer.Text)
	}

	switch p.PollType {
	case poll.SingleChoice:
		options = append(options, back)
		idx, _, _ := selectOption("Select answer", options)
		if idx >= len(answers) {
			return
		}

		votes = append(votes, answers[idx].ID)
	case poll.MultipleChoice:
		options = append(options, sendVotes)
		options = append(options, back)

		done := false
		for !done {
			idx, _, _ := selectOption("Select answer", options)
			if idx < len(answers) {
				votes = append(votes, answers[idx].ID)
				options = append(options[:idx], options[idx+1:]...)
				answers = append(answers[:idx], answers[idx+1:]...)
			} else if idx == len(answers) {
				done = true
			} else if idx > len(answers) {
				return
			}
		}
	}

	err := pollService.Vote(p.ID, &votes, voterID)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	clearScreen()
}

func getResults(p *poll.Poll) {
	results, err := pollService.GetResults(p.ID)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	logPollDetails(p)
	for _, result := range *results {
		fmt.Println(fmt.Sprintf("%s - %d", result.Answer.Text, result.VotesNo))
	}
}

func logPollDetails(p *poll.Poll) {
	fmt.Println("----------------")
	fmt.Println(fmt.Sprintf("Title: %s", p.Title))
	fmt.Println(fmt.Sprintf("Description: %s", p.Description))
	fmt.Println(fmt.Sprintf("DueTime: %s", p.DueTime))
	fmt.Println(fmt.Sprintf("PollType: %s", pollTypeToString[p.PollType]))
	fmt.Println("----------------")
}

func main() {
	addresses := strings.Split(os.Getenv("ADDRESS"), ",")
	if len(addresses) == 1 && addresses[0] == "" {
		log.Fatal("ADDRESS env variable not specified")
		os.Exit(1)
	}
	keyspace := os.Getenv("KEYSPACE")
	if keyspace == "" {
		log.Fatal("KEYSPACE env variable not specified")
		os.Exit(1)
	}

	repo, _ := cass.NewCassandraRepo(addresses, keyspace)
	pollService = poll.NewPollService(repo)

	voterID, _ = gocql.RandomUUID()

	clearScreen()

	for {
		menu()
	}
}
