package main

import (
	"fmt"
	"time"

	"github.com/gocql/gocql"
	"github.com/roertbb/casspoll/poll"
	cass "github.com/roertbb/casspoll/repo/cassandra"
)

func main() {
	repo, _ := cass.NewCassandraRepo()
	service := poll.NewPollService(repo)

	pollID, _ := gocql.RandomUUID()
	examplePoll := poll.Poll{
		ID:          pollID,
		Title:       "test",
		Description: "test desc",
		PollType:    poll.SingleChoice,
		DueTime:     time.Now().Add(time.Second * 15),
	}

	answers := []poll.Answer{}
	for _, idx := range []string{"1", "2", "3"} {
		id, _ := gocql.RandomUUID()
		answers = append(answers, poll.Answer{ID: id, Text: idx, PollID: pollID})
	}

	service.CreatePoll(&examplePoll, &answers)

	polls, _ := service.GetActivePolls()
	fmt.Println("---")
	fmt.Println("active polls")
	fmt.Println(polls)

	pollAnswers, _ := service.GetAnswers(pollID)
	fmt.Println("---")
	fmt.Println("answers")
	fmt.Println(pollAnswers)

	voterID, _ := gocql.RandomUUID()
	service.Vote(&poll.Vote{AnswerID: answers[0].ID, PollID: examplePoll.ID, VoterID: voterID})
	service.Vote(&poll.Vote{AnswerID: answers[1].ID, PollID: examplePoll.ID, VoterID: voterID})

	voter2ID, _ := gocql.RandomUUID()
	service.Vote(&poll.Vote{AnswerID: answers[0].ID, PollID: examplePoll.ID, VoterID: voter2ID})

	results, _ := service.GetResults(examplePoll.ID, examplePoll.DueTime)

	fmt.Println("---")
	fmt.Println("results")
	fmt.Println(results)
}
