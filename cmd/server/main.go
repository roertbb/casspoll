package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"github.com/roertbb/casspoll/poll"
	cass "github.com/roertbb/casspoll/repo/cassandra"
)

var pollService poll.PollService
var voterID gocql.UUID

func getActivePolls(w http.ResponseWriter, r *http.Request) {
	polls, err := pollService.GetActivePolls()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(polls)
}

func createPoll(w http.ResponseWriter, r *http.Request) {
	type createPollDTO struct {
		Title       string        `json:"title"`
		Description string        `json:"description"`
		PollType    poll.PollType `json:"type"`
		DueTime     time.Time     `json:"dueTime"`
		Answers     []string      `json:"answers"`
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var data createPollDTO
	err := json.Unmarshal(reqBody, &data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to parse request body"})
		return
	}

	p := poll.Poll{
		Title:       data.Title,
		Description: data.Description,
		PollType:    data.PollType,
		DueTime:     data.DueTime,
	}

	uuid, err := pollService.CreatePoll(&p, &data.Answers)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]gocql.UUID{"id": uuid})
}

func getPollByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid, err := gocql.ParseUUID(vars["uuid"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to find poll with given id"})
		return
	}

	poll, err := pollService.GetPollByID(uuid)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(poll)
}

func getAnswers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid, err := gocql.ParseUUID(vars["uuid"])

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to find poll with given id"})
		return
	}

	answers, err := pollService.GetAnswers(uuid)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(answers)
}

func getResults(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid, err := gocql.ParseUUID(vars["uuid"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to find poll with given id"})
		return
	}

	results, err := pollService.GetResults(uuid)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results)
}

func vote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid, err := gocql.ParseUUID(vars["uuid"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to find poll with given id"})
		return
	}

	type voteDTO struct {
		VoterID string   `json:"voterId"`
		Answers []string `json:"answers"`
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var data voteDTO
	err = json.Unmarshal(reqBody, &data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to parse request body"})
		return
	}

	voterUUID, _ := gocql.ParseUUID(data.VoterID)
	answerUUIDs := []gocql.UUID{}
	for _, ans := range data.Answers {
		ansUUID, _ := gocql.ParseUUID(ans)
		answerUUIDs = append(answerUUIDs, ansUUID)
	}

	err = pollService.Vote(uuid, &answerUUIDs, voterUUID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	addresses := strings.Split(os.Getenv("ADDRESS"), ",")
	if len(addresses) == 1 && addresses[0] == "" {
		log.Fatal("ADDRESS env variable not specified")
		os.Exit(1)
	}
	keyspace := os.Getenv("KEYSPACE")
	if keyspace == "" {
		log.Fatal("KEYSPACE env variable not specified")
		os.Exit(1)
	}

	repo, _ := cass.NewCassandraRepo(addresses, keyspace)
	pollService = poll.NewPollService(repo)

	r := mux.NewRouter()
	r.HandleFunc("/polls", getActivePolls).Methods("GET")
	r.HandleFunc("/polls", createPoll).Methods("POST")
	r.HandleFunc("/polls/{uuid}", getPollByID).Methods("GET")
	r.HandleFunc("/polls/{uuid}/answers", getAnswers).Methods("GET")
	r.HandleFunc("/polls/{uuid}/results", getResults).Methods("GET")
	r.HandleFunc("/polls/{uuid}/vote", vote).Methods("POST")

	fmt.Println("server running - 127.0.0.1:8080")

	log.Fatal(http.ListenAndServe(":8080", r))
}
