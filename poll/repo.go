package poll

import (
	"time"

	"github.com/gocql/gocql"
)

type PollRepo interface {
	CreatePoll(poll *Poll) error
	CreateAnswer(answer *Answer) error
	CreateVote(vote *Vote, timestamp time.Time) error
	GetResults(pollID gocql.UUID, dueTime time.Time) (*map[gocql.UUID]int, error)
	// GetPolls() (*[]Poll, error)
	// GetAnswers(poll *Poll) (*[]Answer, error)
}
