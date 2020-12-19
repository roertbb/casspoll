package cassandra

import (
	"log"
	"time"

	"github.com/gocql/gocql"
	"github.com/roertbb/casspoll/poll"
)

type cassandraRepo struct {
	cluster *gocql.ClusterConfig
	session *gocql.Session
}

func NewCassandraRepo() (*cassandraRepo, error) {
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "casspoll"
	cluster.Consistency = gocql.One
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}

	repo := &cassandraRepo{
		cluster: cluster,
		session: session,
	}

	return repo, nil
}

func (c *cassandraRepo) GetActivePolls(timestamp time.Time) (*[]poll.Poll, error) {
	var pollID gocql.UUID
	var title, description string
	var dueTime time.Time
	var pollType poll.PollType

	activePolls := []poll.Poll{}

	iter := c.session.Query(`SELECT pollId, title, description, dueTime, pollType FROM ActivePolls`).Consistency(gocql.One).Iter()
	for iter.Scan(&pollID, &title, &description, &dueTime, &pollType) {
		activePolls = append(activePolls, poll.Poll{
			ID:          pollID,
			Title:       title,
			Description: description,
			PollType:    pollType,
			DueTime:     dueTime,
		})
	}
	if err := iter.Close(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &activePolls, nil
}

func (c *cassandraRepo) CreatePoll(p *poll.Poll) error {
	err := c.session.Query(`INSERT INTO polls (pollId, title, description, dueTime, pollType) VALUES (?, ?, ?, ?, ?)`, p.ID, p.Title, p.Description, p.DueTime, p.PollType).Exec()
	if err != nil {
		return err
	}

	now := time.Now()
	ttl := int(p.DueTime.Sub(now).Seconds())

	err = c.session.Query(`INSERT INTO ActivePolls (pollId, title, description, dueTime, pollType) VALUES (?, ?, ?, ?, ?) USING TTL ?`, p.ID, p.Title, p.Description, p.DueTime, p.PollType, ttl).Exec()

	return err
}

func (c *cassandraRepo) GetAnswersByPollID(pollID gocql.UUID) (*[]poll.Answer, error) {
	var answerID gocql.UUID
	var text string

	answers := []poll.Answer{}

	iter := c.session.Query(`SELECT answerId, answer FROM answers WHERE pollId = ?`, pollID).Consistency(gocql.One).Iter()
	for iter.Scan(&answerID, &text) {
		answers = append(answers, poll.Answer{ID: answerID, Text: text, PollID: pollID})
	}
	if err := iter.Close(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &answers, nil
}

func (c *cassandraRepo) CreateAnswer(answer *poll.Answer) error {
	return c.session.Query(`INSERT INTO answers (answerId, answer, pollId) VALUES (?, ?, ?)`, answer.ID, answer.Text, answer.PollID).Exec()
}

func (c *cassandraRepo) CreateVote(vote *poll.Vote, timestamp time.Time) error {
	return c.session.Query(`INSERT INTO votes (answerId, pollId, voterId) VALUES (?, ?, ?)`, vote.AnswerID, vote.PollID, vote.VoterID).Exec()
}

func (c *cassandraRepo) GetResults(pollID gocql.UUID, dueTime time.Time) (*map[gocql.UUID]int, error) {
	var answerID gocql.UUID
	var votesNo int
	results := map[gocql.UUID]int{}

	// TODO: how to handle check if timestamp < dueTime without ALLOW FILTERING
	// https://www.datastax.com/blog/new-cassandra-30-materialized-views
	// potentially can skip that condition and prevent from posting votes after dueTime is met from inside service layer, quack :v
	// iter := c.session.Query(`SELECT answerId, COUNT(*) FROM votes WHERE pollId = ? AND createdAt < ? GROUP BY answerId ALLOW FILTERING`, pollID, dueTime).Consistency(gocql.One).Iter()

	iter := c.session.Query(`SELECT answerId, COUNT(*) FROM votes WHERE pollId = ? GROUP BY answerId`, pollID).Consistency(gocql.One).Iter()
	for iter.Scan(&answerID, &votesNo) {
		results[answerID] = votesNo
	}
	if err := iter.Close(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &results, nil
}
