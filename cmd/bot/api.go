package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gocql/gocql"
	"github.com/roertbb/casspoll/poll"
)

// TODO: customizable?
const minAnwers = 3
const maxAnswers = 9

func createRandomPoll(address string, dueTime time.Time) (*poll.Poll, error) {
	type createPollDTO struct {
		Title       string        `json:"title"`
		Description string        `json:"description"`
		PollType    poll.PollType `json:"type"`
		DueTime     time.Time     `json:"dueTime"`
		Answers     []string      `json:"answers"`
	}

	rand.Seed(time.Now().UnixNano())

	answersNo := rand.Intn(maxAnswers-minAnwers) + minAnwers
	answers := []string{}
	for i := 0; i < answersNo; i++ {
		answers = append(answers, randString(30))
	}

	data := createPollDTO{
		Title:       randString(20),
		Description: randString(50),
		PollType:    poll.PollType(rand.Intn(2) + 1),
		DueTime:     dueTime,
		Answers:     answers,
	}

	marshBody, _ := json.Marshal(data)

	resp, err := http.Post(address+"/polls", "application/json", bytes.NewBuffer(marshBody))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	scanner.Scan()
	body := scanner.Text()

	type response struct {
		ID gocql.UUID `json:"id"`
	}

	responseData := response{}
	err = json.Unmarshal([]byte(body), &responseData)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	p := poll.Poll{
		ID:          responseData.ID,
		Title:       data.Title,
		Description: data.Description,
		PollType:    data.PollType,
		DueTime:     data.DueTime,
	}

	return &p, nil
}

func getActivePolls(address string) (*[]poll.Poll, error) {
	resp, err := http.Get(address + "/polls")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)
	scanner.Scan()
	respBody := scanner.Text()

	pollsData := []poll.Poll{}
	err = json.Unmarshal([]byte(respBody), &pollsData)

	return &pollsData, nil
}

func getAnswers(address string, pollData *poll.Poll) (*[]poll.Answer, error) {
	resp, err := http.Get(address + "/polls/" + pollData.ID.String() + "/answers")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)
	scanner.Scan()
	respBody := scanner.Text()

	answersData := []poll.Answer{}
	err = json.Unmarshal([]byte(respBody), &answersData)

	return &answersData, nil
}

func voteRandomlyForPoll(address string, pollData *poll.Poll, voterID gocql.UUID) (*[]gocql.UUID, error) {
	type voteDTO struct {
		VoterID gocql.UUID   `json:"voterId"`
		Answers []gocql.UUID `json:"answers"`
	}

	rand.Seed(time.Now().UnixNano())

	resp, err := http.Get(address + "/polls/" + pollData.ID.String() + "/answers")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)
	scanner.Scan()
	respBody := scanner.Text()

	answersData := []poll.Answer{}
	err = json.Unmarshal([]byte(respBody), &answersData)

	selectedAnswers := []gocql.UUID{}
	if pollData.PollType == poll.MultipleChoice {
		randID := rand.Intn(len(answersData))
		selectedAnswers = append(selectedAnswers, answersData[randID].ID)
	} else {
		IDs := rand.Perm(len(answersData))
		for _, idx := range IDs {
			selectedAnswers = append(selectedAnswers, answersData[idx].ID)
		}
	}

	body := voteDTO{VoterID: voterID, Answers: selectedAnswers}
	marshBody, err := json.Marshal(body)

	resp, err = http.Post(address+"/polls/"+pollData.ID.String()+"/vote", "application/json", bytes.NewBuffer(marshBody))

	if resp.StatusCode != http.StatusOK {
		scanner := bufio.NewScanner(resp.Body)
		scanner.Scan()
		respBody := scanner.Text()
		responseError := map[string]string{}
		err = json.Unmarshal([]byte(respBody), &responseError)
		errorMessage := responseError["error"]
		fmt.Println(errorMessage)
		return nil, errors.New(errorMessage)
	}

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &selectedAnswers, nil
}

func getPollResults(address string, pollData *poll.Poll) (*[]poll.Result, error) {
	resp, err := http.Get(address + "/polls/" + pollData.ID.String() + "/results")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)
	scanner.Scan()
	respBody := scanner.Text()

	resultsData := []poll.Result{}
	err = json.Unmarshal([]byte(respBody), &resultsData)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &resultsData, nil
}
