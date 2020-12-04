package poll

type PollRepo interface {
	CreatePoll(poll *Poll, answers *[]Answer) error
	CreateVote(poll *Poll, answers *[]Answer) error
	GetResults(poll *Poll) (*[]Result, error)
	// GetPolls() (*[]Poll, error)
	// GetAnswers(poll *Poll) (*[]Answer, error)
}
