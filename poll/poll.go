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

func (p *pollService) GetActivePolls() (*[]Poll, error) {
	return p.repo.GetActivePolls()
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

func (p *pollService) GetAnswers(pollID gocql.UUID) (*[]Answer, error) {
	return p.repo.GetAnswersByPollID(pollID)
}

func (p *pollService) Vote(vote *Vote) error {
	return p.repo.CreateVote(vote, time.Now())
}

func (p *pollService) GetResults(pollID gocql.UUID, dueTime time.Time) (*map[Answer]int, error) {
	answers, err := p.repo.GetAnswersByPollID(pollID)
	if err != nil {
		return nil, err
	}

	results, _ := p.repo.GetResults(pollID, dueTime)
	if err != nil {
		return nil, err
	}

	pollResults := map[Answer]int{}
	for _, answer := range *answers {
		pollResults[answer] = (*results)[answer.ID]
	}

	return &pollResults, nil
}
