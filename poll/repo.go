package poll

import "github.com/gocql/gocql"

type PollRepo interface {
	CreatePoll(poll *Poll) error
	CreateAnswer(answer *Answer, pollID gocql.UUID) error
	CreateVote(pollID, answerID gocql.UUID) error
	GetResults(pollID gocql.UUID) (*[]Result, error)
	// GetPolls() (*[]Poll, error)
	// GetAnswers(poll *Poll) (*[]Answer, error)
}
