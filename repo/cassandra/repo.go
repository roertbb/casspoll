package cassandra

import (
	"fmt"
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

func (c *cassandraRepo) CreatePoll(p *poll.Poll) error {
	return c.session.Query(`INSERT INTO polls (pollId, title, description, dueTime, pollType) VALUES (?, ?, ?, ?, ?)`, p.ID, p.Title, p.Description, p.DueTime, p.PollType).Exec()
}

func (c *cassandraRepo) CreateAnswer(answer *poll.Answer) error {
	return c.session.Query(`INSERT INTO answers (answerId, answer, pollId) VALUES (?, ?, ?)`, answer.ID, answer.Text, answer.PollID).Exec()
}

func (c *cassandraRepo) CreateVote(vote *poll.Vote, timestamp time.Time) error {
	return c.session.Query(`INSERT INTO votes (answerId, pollId, createdAt, voterId) VALUES (?, ?, ?, ?)`, vote.AnswerID, vote.PollID, timestamp, vote.VoterID).Exec()
}

func (c *cassandraRepo) GetResults(pollID gocql.UUID, dueTime time.Time) (*map[gocql.UUID]int, error) {
	var answerID gocql.UUID
	var votesNo int
	results := map[gocql.UUID]int{}

	// TODO: add check if timestamp < dueTime
	// SELECT COUNT(*) FROM Votes WHERE pollId = ? AND timestamp < dueTime GROUP BY answerId
	// https://www.datastax.com/blog/new-cassandra-30-materialized-views

	iter := c.session.Query(`SELECT answerId, COUNT(voterId) FROM votes WHERE pollId = ? GROUP BY answerId`, pollID).Consistency(gocql.One).Iter()
	for iter.Scan(&answerID, &votesNo) {
		fmt.Println(answerID, votesNo)
		results[answerID] = votesNo
	}
	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}

	return &results, nil
}
