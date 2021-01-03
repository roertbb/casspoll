package poll

import (
	"time"

	"github.com/gocql/gocql"
)

type PollType int

const (
	SingleChoice PollType = iota
	MultipleChoice
)

type Poll struct {
	ID          gocql.UUID `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	PollType    PollType   `json:"type"`
	DueTime     time.Time  `json:"dueTime"`
}

type Answer struct {
	ID     gocql.UUID `json:"id"`
	Text   string     `json:"text"`
	PollID gocql.UUID `json:"-"`
}

type Vote struct {
	AnswerID gocql.UUID `json:"answerId"`
	VoterID  gocql.UUID `json:"voterId"`
	PollID   gocql.UUID `json:"-"`
}

type Result struct {
	Answer  *Answer `json:"answers"`
	VotesNo int     `json:"votesNo"`
}
