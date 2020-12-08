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

	id, _ := gocql.RandomUUID()
	examplePoll := poll.Poll{
		ID:          id,
		Title:       "test",
		Description: "test desc",
		PollType:    poll.SingleChoice,
		DueTime:     time.Now(),
	}

	answers := []poll.Answer{}
	for _, idx := range []string{"1", "2", "3"} {
		id, _ := gocql.RandomUUID()
		answers = append(answers, poll.Answer{ID: id, Text: idx})
	}

	service.CreatePoll(&examplePoll, &answers)

	service.Vote(examplePoll.ID, answers[0].ID)
	service.Vote(examplePoll.ID, answers[0].ID)
	service.Vote(examplePoll.ID, answers[1].ID)

	results, _ := service.GetResults(examplePoll.ID)
	fmt.Println(results)
}
