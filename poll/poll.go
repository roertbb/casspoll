package poll

import (
	"log"
	"time"

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
		err := p.repo.CreateAnswer(&answer)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}

	return nil
}

func (p *pollService) Vote(vote *Vote) error {
	timestamp := time.Now()
	return p.repo.CreateVote(vote, timestamp)
}

func (p *pollService) GetResults(pollID gocql.UUID, dueTime time.Time) (*map[gocql.UUID]int, error) {
	return p.repo.GetResults(pollID, dueTime)
}
