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
		PollType:    poll.MultipleChoice,
		DueTime:     time.Now().Add(time.Second * 10),
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

	votes := []poll.Vote{
		{AnswerID: answers[0].ID, PollID: examplePoll.ID, VoterID: voterID},
		{AnswerID: answers[1].ID, PollID: examplePoll.ID, VoterID: voterID},
	}
	service.Vote(&examplePoll, &votes)

	voter2ID, _ := gocql.RandomUUID()
	votes = []poll.Vote{
		{AnswerID: answers[0].ID, PollID: examplePoll.ID, VoterID: voter2ID},
	}
	service.Vote(&examplePoll, &votes)

	results, _ := service.GetResults(examplePoll.ID, examplePoll.DueTime)

	fmt.Println("---")
	fmt.Println("results")
	fmt.Println(results)
}
