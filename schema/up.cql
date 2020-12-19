CREATE TABLE IF NOT EXISTS Polls (
    pollId uuid,
    title text,
    description text,
    pollType int,
    dueTime timestamp,

    PRIMARY KEY (pollId, dueTime)
) WITH CLUSTERING ORDER BY (dueTime DESC);

CREATE TABLE IF NOT EXISTS ActivePolls (
    pollId uuid,
    title text,
    description text,
    pollType int,
    dueTime timestamp,

    PRIMARY KEY (pollId, dueTime)
) WITH CLUSTERING ORDER BY (dueTime DESC);

CREATE TABLE IF NOT EXISTS Answers (
    answerId uuid,
    answer text,
    pollId uuid,

    PRIMARY KEY (pollId, answerId)
);

CREATE TABLE IF NOT EXISTS Votes (
    pollId uuid,
    answerId uuid,
    voterId uuid,
    
    PRIMARY KEY (pollId, answerId, voterId)
);

-- TODO: Votes by voterId