package poll

type PollType int

const (
	singleChoice PollType = iota
	multipleChoice
)

type Poll struct {
	ID          int
	Title       string
	Description string
	PollType    PollType
	DueTime     int
}

type Answer struct {
	ID   int
	Text string
}

type Result struct {
	Answer  Answer
	VotesNo int
}
