package poll

type PollService interface {
	CreatePoll(poll *Poll, answers *[]Answer) error
	Vote(poll *Poll, answers *[]Answer) error
	GetResults(poll *Poll) (*[]Result, error)
	// ListPolls() (*[]Poll, error)
	// ListAnswers(poll *Poll) (*[]Answer, error)
}
