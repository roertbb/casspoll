package poll

import (
	"errors"
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

func (p *pollService) Vote(pollData *Poll, votes *[]Vote) error {
	now := time.Now()
	if now.After(pollData.DueTime) {
		return errors.New("Cannot vote in the poll after it's due time")
	}

	if len(*votes) == 0 {
		return errors.New("You need to select at least one answer")
	}

	if pollData.PollType == SingleChoice && len(*votes) > 1 {
		return errors.New("Cannot select multiple answers for single-choice poll")
	}

	for _, vote := range *votes {
		err := p.repo.CreateVote(&vote, now)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}

	return nil
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
