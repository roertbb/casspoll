package poll

import (
	"errors"
	"fmt"
	"sort"
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

func (p *pollService) GetPollByID(pollID gocql.UUID) (*Poll, error) {
	return p.repo.GetPollByID(pollID)
}

func (p *pollService) CreatePoll(poll *Poll, answers *[]string) (gocql.UUID, error) {
	pollID, _ := gocql.RandomUUID()
	poll.ID = pollID

	if time.Now().After(poll.DueTime) {
		return gocql.UUID{}, errors.New("Cannot create poll with due time with the past")
	}

	err := p.repo.CreatePoll(poll)
	if err != nil {
		fmt.Println(err)
		return gocql.UUID{}, err
	}

	for _, answer := range *answers {
		answerID, _ := gocql.RandomUUID()
		a := Answer{ID: answerID, Text: answer, PollID: pollID}
		err := p.repo.CreateAnswer(&a)
		if err != nil {
			fmt.Println(err)
			return gocql.UUID{}, err
		}
	}

	return pollID, nil
}

func (p *pollService) GetAnswers(pollID gocql.UUID) (*[]Answer, error) {
	return p.repo.GetAnswersByPollID(pollID)
}

func (p *pollService) Vote(pollID gocql.UUID, answerIDs *[]gocql.UUID, voterID gocql.UUID) error {
	now := time.Now()

	pollData, err := p.repo.GetPollByID(pollID)
	if err != nil {
		return errors.New("There is no poll with given id")
	}

	if now.After(pollData.DueTime) {
		return errors.New("Cannot vote in the poll after it's due time")
	}

	if len(*answerIDs) == 0 {
		return errors.New("You need to select at least one answer")
	}

	if pollData.PollType == SingleChoice && len(*answerIDs) > 1 {
		return errors.New("Cannot select multiple answers for single-choice poll")
	}

	for idx := range *answerIDs {
		vote := Vote{
			PollID:   pollID,
			AnswerID: (*answerIDs)[idx],
			VoterID:  voterID,
		}

		err := p.repo.CreateVote(&vote, now)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}

func (p *pollService) GetResults(pollID gocql.UUID) (*[]Result, error) {
	answers, err := p.repo.GetAnswersByPollID(pollID)
	if err != nil {
		return nil, err
	}

	results, _ := p.repo.GetResults(pollID)
	if err != nil {
		return nil, err
	}

	resSlice := []Result{}
	for idx := range *answers {
		curAnswer := (*answers)[idx]
		resSlice = append(resSlice, Result{Answer: &curAnswer, VotesNo: (*results)[curAnswer.ID]})
	}

	sort.Slice(resSlice, func(i, j int) bool {
		return resSlice[i].VotesNo > resSlice[j].VotesNo
	})

	return &resSlice, nil
}
