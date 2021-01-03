package poll

import (
	"time"

	"github.com/gocql/gocql"
)

type PollRepo interface {
	GetActivePolls() (*[]Poll, error)
	CreatePoll(poll *Poll) error
	GetPollByID(pollID gocql.UUID) (*Poll, error)
	GetAnswersByPollID(pollID gocql.UUID) (*[]Answer, error)
	CreateAnswer(answer *Answer) error
	CreateVote(vote *Vote, timestamp time.Time) error
	GetResults(pollID gocql.UUID) (*map[gocql.UUID]int, error)
}
