package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/roertbb/casspoll/poll"
)

func main() {
	addresses := strings.Split(os.Getenv("ADDRESS"), ",")
	if len(addresses) == 1 && addresses[0] == "" {
		log.Fatal("ADDRESS env variable not specified")
		os.Exit(1)
	}

	// addresses := []string{"http://127.0.0.1:8080"}

	// config
	pollsNum := 40
	partitionStartsIn := 15
	pollEndInSeconds := 30
	syncAfterSeconds := 20
	// voterNo := 1

	votingDone := false

	createChan := make(chan bool)
	answersChan := make(chan bool)
	doneChan := make(chan bool)

	// create 100 polls with dueTime in 2 min
	dueTime := time.Now().Add(time.Second * time.Duration(pollEndInSeconds))

	for i := 0; i < pollsNum; i++ {
		go func() {
			createRandomPoll(addresses[0], dueTime)
			createChan <- true
		}()
	}

	for i := 0; i < pollsNum; i++ {
		<-createChan
	}

	polls, _ := getActivePolls(addresses[0])
	realResults := make(map[gocql.UUID](map[gocql.UUID]int))
	for _, poll := range *polls {
		realResults[poll.ID] = make(map[gocql.UUID]int)
	}

	for idx := range *polls {
		poll := (*polls)[idx]

		go func() {
			answers, _ := getAnswers(addresses[0], &poll)
			for aidx := range *answers {
				ans := (*answers)[aidx]
				realResults[poll.ID][ans.ID] = 0
			}
			answersChan <- true
		}()
	}

	for i := 0; i < pollsNum; i++ {
		<-answersChan
	}

	startPartitionTimer := time.NewTimer(time.Second * time.Duration(partitionStartsIn))
	go func() {
		<-startPartitionTimer.C
		fmt.Println("partition start")
		partitionCmd := exec.Command("docker", "exec", "cass1", "/bin/bash", "-c", "\"/utils/start-partition.sh\"")
		partitionCmd.Start()
		partitionCmd.Wait()
	}()

	finishVotingTimer := time.NewTimer(time.Second * time.Duration(pollEndInSeconds))
	go func() {
		<-finishVotingTimer.C
		fmt.Println("finish voting")
		votingDone = true

		res, _ := getPollsResults(addresses, polls)
		wrongAnswersCount, wrongPollWinner := resultsSummary(&realResults, res)

		fmt.Println("after voting time finished")
		fmt.Println("wrongAnswersCount", wrongAnswersCount)
		fmt.Println("wrongPollWinner", wrongPollWinner)

		partitionEnd := exec.Command("docker", "exec", "cass1", "/bin/bash", "-c", "\"/utils/stop-partition.sh\"")
		partitionEnd.Start()
		partitionEnd.Wait()

		letItSyncTimer := time.NewTimer(time.Second * time.Duration(syncAfterSeconds))
		<-letItSyncTimer.C

		res, _ = getPollsResults(addresses, polls)
		wrongAnswersCount, wrongPollWinner = resultsSummary(&realResults, res)

		fmt.Println("after removing partition and syncing")
		fmt.Println("wrongAnswersCount", wrongAnswersCount)
		fmt.Println("wrongPollWinner", wrongPollWinner)

		doneChan <- true
	}()

	// TODO: spin more than one thread and gather answers from them
	// spin up 10 voting threads - save locally what they voted for
	// for i := 0; i < voterNo; i++ {
	go func() {
		for !votingDone {
			addressID := rand.Intn(len(addresses))
			pollID := rand.Intn(len(*polls))
			voterUUID, _ := gocql.RandomUUID()

			currentPoll := (*polls)[pollID]
			selectedVotes, err := voteRandomlyForPoll(addresses[addressID], &currentPoll, voterUUID)
			if err == nil {
				for _, answerID := range *selectedVotes {
					realResults[currentPoll.ID][answerID]++
				}
			}

		}
	}()
	// }

	<-doneChan
}

func getPollsResults(addresses []string, polls *[]poll.Poll) (*map[gocql.UUID](map[gocql.UUID]int), error) {
	answerCount := make(map[gocql.UUID](map[gocql.UUID]int))
	for _, poll := range *polls {
		answerCount[poll.ID] = make(map[gocql.UUID]int)
	}

	for _, address := range addresses {
		for _, poll := range *polls {
			result, _ := getPollResults(address, &poll)

			for _, r := range *result {
				answerCount[poll.ID][r.Answer.ID] = r.VotesNo
			}
		}
	}

	return &answerCount, nil
}

func resultsSummary(realResults *map[gocql.UUID](map[gocql.UUID]int), results *map[gocql.UUID](map[gocql.UUID]int)) (int, int) {
	wrongAnswersCount := 0
	wrongPollWinner := 0

	for pollID, answers := range *realResults {

		realResultsWinner := gocql.UUID{}
		realResultsWinnerVotesNo := 0

		resultsWinner := gocql.UUID{}
		resultsWinnerVotesNo := 0

		for answersID, realVotesNo := range *&answers {
			votesNo := (*results)[pollID][answersID]

			if votesNo != realVotesNo {
				wrongAnswersCount++

				fmt.Printf("PollID: %s | AnswerID: %s | wrong answer difference is - %d / expected - %d\n", pollID.String(), answersID.String(), votesNo, realVotesNo)

				if realVotesNo > realResultsWinnerVotesNo {
					realResultsWinner = answersID
					realResultsWinnerVotesNo = realVotesNo
				}

				if votesNo > resultsWinnerVotesNo {
					resultsWinner = answersID
					resultsWinnerVotesNo = votesNo
				}
			}
		}

		if realResultsWinner != resultsWinner {
			wrongPollWinner++

			fmt.Printf("PollID: %s | wrong poll winner | is - %s / expected - %s\n", pollID.String(), resultsWinner, realResultsWinner)
		}
	}

	return wrongAnswersCount, wrongPollWinner
}
