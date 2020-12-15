package poll

import (
	"time"

	"github.com/gocql/gocql"
)

type PollRepo interface {
	GetActivePolls(timestamp time.Time) (*[]Poll, error)
	CreatePoll(poll *Poll) error
	GetAnswersByPollID(pollID gocql.UUID) (*[]Answer, error)
	CreateAnswer(answer *Answer) error
	CreateVote(vote *Vote, timestamp time.Time) error
	GetResults(pollID gocql.UUID, dueTime time.Time) (*map[gocql.UUID]int, error)
}
