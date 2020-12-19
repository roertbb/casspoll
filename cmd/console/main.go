package main

import (
	"fmt"
	"log"
	"os"

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
)

var pollService poll.PollService
var pollTypeToString map[poll.PollType]string = map[poll.PollType]string{
	0: singleChoice,
	1: multipleChoice,
}

func menu() {
	result, _ := selectStringOption("Select option", []string{
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

	pollID, _ := gocql.RandomUUID()
	newPoll := poll.Poll{
		ID: pollID,
	}

	if newPoll.Title, err = promptNonEmptyString("Enter title"); err != nil {
		log.Fatal(err)
		return
	}

	if newPoll.Description, err = promptNonEmptyString("Enter description"); err != nil {
		log.Fatal(err)
		return
	}

	if newPoll.DueTime, err = promptFutureDate("Enter due time (in YYYY-MM-DDTHH:MM:SS format)"); err != nil {
		log.Fatal(err)
		return
	}

	prompt := promptui.Select{
		Label: "Select type",
		Items: []string{singleChoice, multipleChoice},
	}

	choiceType, _, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
	newPoll.PollType = poll.PollType(choiceType)

	fmt.Println("----------------")
	fmt.Println("Enter answers for the poll")
	fmt.Println("(type and confirm with <Enter> / leave empty and press <Enter> if all answers are added)")
	fmt.Println("----------------")

	answers := []poll.Answer{}
	done := false
	for !done {
		answer, err := promptString("Enter answer")
		if err != nil {
			log.Fatal(err)
			return
		}

		if answer != "" {
			answerID, _ := gocql.RandomUUID()
			answers = append(answers, poll.Answer{ID: answerID, Text: answer, PollID: newPoll.ID})
		} else {
			done = true
		}
	}

	reverse(answers)

	err = pollService.CreatePoll(&newPoll, &answers)
	if err != nil {
		log.Fatal(err)
		return
	}

	clearScreen()

	fmt.Println("Successfully created poll!")

	menu()
}

func listActivePolls() {
	polls, err := pollService.GetActivePolls()
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
	pollList = append(pollList, back)

	result, _ := selectStringOption("Select poll", pollList)

	clearScreen()

	switch result {
	case back:
		menu()
	default:
		selectedPoll := pollMap[result]
		pollDetails(&selectedPoll)
	}
}

func pollDetails(p *poll.Poll) {
	answers, err := pollService.GetAnswers(p.ID)
	if err != nil {
		log.Fatal(err)
		return
	}

	logPollDetails(p)
	fmt.Println("Answers:")
	for _, answer := range *answers {
		fmt.Println(fmt.Sprintf("- %s", answer.Text))
	}

	result, _ := selectStringOption("Select option", []string{votePoll, getPollResults, back})

	clearScreen()

	switch result {
	case back:
		listActivePolls()
	case votePoll:
		vote(p, answers)
	case getPollResults:
		getResults(p)
	}
}

func vote(p *poll.Poll, answers *[]poll.Answer) {
	switch p.PollType {
	case poll.SingleChoice:
		fmt.Println("TODO: poll.SingleChoice")
	case poll.MultipleChoice:
		fmt.Println("TODO: poll.MultipleChoice")

	}
}

func getResults(p *poll.Poll) {
	fmt.Println("TODO: get poll results")
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
	repo, _ := cass.NewCassandraRepo()
	pollService = poll.NewPollService(repo)

	clearScreen()
	menu()
}
