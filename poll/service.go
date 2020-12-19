package poll

import (
	"time"

	"github.com/gocql/gocql"
)

type PollService interface {
	GetActivePolls() (*[]Poll, error)
	CreatePoll(poll *Poll, answers *[]Answer) error
	GetAnswers(pollID gocql.UUID) (*[]Answer, error)
	Vote(poll *Poll, votes *[]Vote) error
	GetResults(pollID gocql.UUID, dueTime time.Time) (*map[Answer]int, error)
}
