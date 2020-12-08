package poll

import "github.com/gocql/gocql"

type PollService interface {
	CreatePoll(poll *Poll, answers *[]Answer) error
	Vote(pollID, answerID gocql.UUID) error
	GetResults(pollID gocql.UUID) (*[]Result, error)
	// ListPolls() (*[]Poll, error)
	// ListAnswers(poll *Poll) (*[]Answer, error)
}
