package poll

import (
	"github.com/gocql/gocql"
)

type PollService interface {
	GetActivePolls() (*[]Poll, error)
	CreatePoll(poll *Poll, answers *[]Answer) (gocql.UUID, error)
	GetAnswers(pollID gocql.UUID) (*[]Answer, error)
	Vote(pollID gocql.UUID, answerIDs *[]gocql.UUID, voterID gocql.UUID) error
	GetResults(pollID gocql.UUID) (*[]Result, error)
}
