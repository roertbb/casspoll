package poll

import (
	"log"

	"github.com/gocql/gocql"
)

type pollService struct {
	repo PollRepo
}

func NewPollService(pollRepo PollRepo) PollService {
	return &pollService{
		pollRepo,
	}
}

func (p *pollService) CreatePoll(poll *Poll, answers *[]Answer) error {
	err := p.repo.CreatePoll(poll)
	if err != nil {
		log.Fatal(err)
		return err
	}

	for _, answer := range *answers {
		err := p.repo.CreateAnswer(&answer, poll.ID)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}

	return nil
}

func (p *pollService) Vote(pollID, answerID gocql.UUID) error {
	return p.repo.CreateVote(pollID, answerID)
}

func (p *pollService) GetResults(pollID gocql.UUID) (*[]Result, error) {
	return p.repo.GetResults(pollID)
}
