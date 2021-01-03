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

func NewCassandraRepo(addresses []string, keyspace string) (*cassandraRepo, error) {
	cluster := gocql.NewCluster(addresses...)
	cluster.Keyspace = keyspace
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

func (c *cassandraRepo) GetActivePolls() (*[]poll.Poll, error) {
	var pollID gocql.UUID
	var title, description string
	var dueTime time.Time
	var pollType poll.PollType

	activePolls := []poll.Poll{}

	iter := c.session.Query(`SELECT poll_id, title, description, due_time, poll_type FROM active_polls`).Consistency(gocql.One).Iter()
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
	err := c.session.Query(`INSERT INTO polls (poll_id, title, description, due_time, poll_type) VALUES (?, ?, ?, ?, ?)`, p.ID, p.Title, p.Description, p.DueTime, p.PollType).Exec()
	if err != nil {
		return err
	}

	now := time.Now()
	ttl := int(p.DueTime.Sub(now).Seconds())
	timestamp := now.Unix() * 1000

	err = c.session.Query(`INSERT INTO active_polls (poll_id, title, description, due_time, poll_type) VALUES (?, ?, ?, ?, ?) USING TTL ? AND TIMESTAMP ?`, p.ID, p.Title, p.Description, p.DueTime, p.PollType, ttl, timestamp).Exec()

	return err
}

func (c *cassandraRepo) GetAnswersByPollID(pollID gocql.UUID) (*[]poll.Answer, error) {
	var answerID gocql.UUID
	var text string

	answers := []poll.Answer{}

	iter := c.session.Query(`SELECT answer_id, answer FROM answers WHERE poll_id = ?`, pollID).Consistency(gocql.One).Iter()
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
	return c.session.Query(`INSERT INTO answers (answer_id, answer, poll_id) VALUES (?, ?, ?)`, answer.ID, answer.Text, answer.PollID).Exec()
}

func (c *cassandraRepo) CreateVote(vote *poll.Vote, now time.Time) error {
	timestamp := now.Unix() * 1000
	return c.session.Query(`INSERT INTO votes (answer_id, poll_id, voter_id) VALUES (?, ?, ?) USING TIMESTAMP ?`, vote.AnswerID, vote.PollID, vote.VoterID, timestamp).Exec()
}

func (c *cassandraRepo) GetResults(pollID gocql.UUID) (*map[gocql.UUID]int, error) {
	var answerID gocql.UUID
	var votesNo int
	results := map[gocql.UUID]int{}

	iter := c.session.Query(`SELECT answer_id, COUNT(*) FROM votes WHERE poll_id = ? GROUP BY answer_id`, pollID).Consistency(gocql.One).Iter()
	for iter.Scan(&answerID, &votesNo) {
		results[answerID] = votesNo
	}
	if err := iter.Close(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &results, nil
}
