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
	ID          gocql.UUID
	Title       string
	Description string
	PollType    PollType
	DueTime     time.Time
}

type Answer struct {
	ID     gocql.UUID
	Text   string
	PollID gocql.UUID
}

type Vote struct {
	PollID   gocql.UUID
	AnswerID gocql.UUID
	VoterID  gocql.UUID
}

type Result struct {
	Answer  *Answer
	VotesNo int
}
