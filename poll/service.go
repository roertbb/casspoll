package poll

import (
	"time"

	"github.com/gocql/gocql"
)

type PollService interface {
	CreatePoll(poll *Poll, answers *[]Answer) error
	Vote(vote *Vote) error
	GetResults(pollID gocql.UUID, dueTime time.Time) (*map[gocql.UUID]int, error)
	// ListPolls() (*[]Poll, error)
	// ListAnswers(poll *Poll) (*[]Answer, error)
}
