package cassandra

import (
	"log"

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

func (c *cassandraRepo) CreatePoll(p *poll.Poll) error {
	return c.session.Query(`INSERT INTO polls (pollId, title, description, dueTime, pollType) VALUES (?, ?, ?, ?, ?)`, p.ID, p.Title, p.Description, p.DueTime, p.PollType).Exec()
}

func (c *cassandraRepo) CreateAnswer(answer *poll.Answer, pollID gocql.UUID) error {
	return c.session.Query(`INSERT INTO answers (answerId, answer, pollId) VALUES (?, ?, ?)`, answer.ID, answer.Text, pollID).Exec()
}

func (c *cassandraRepo) CreateVote(pollID, answerID gocql.UUID) error {
	return c.session.Query(`UPDATE votes SET votesNo = votesNo + 1 WHERE pollId=? AND answerId=?`, pollID, answerID).Exec()
}

func (c *cassandraRepo) GetResults(pollID gocql.UUID) (*[]poll.Result, error) {
	var answerID gocql.UUID
	var votesNo int
	results := []poll.Result{}

	iter := c.session.Query(`SELECT answerId, votesNo FROM votes WHERE pollId=?`, pollID).Consistency(gocql.One).Iter()
	for iter.Scan(&answerID, &votesNo) {
		// TODO: get answers before merge here?
		results = append(results, poll.Result{Answer: poll.Answer{ID: answerID, Text: ""}, VotesNo: votesNo})
	}
	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}

	return &results, nil
}
