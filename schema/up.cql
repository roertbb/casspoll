CREATE TABLE IF NOT EXISTS Polls (
    pollId uuid,
    title text,
    description text,
    dueTime timestamp,
    pollType int,

    PRIMARY KEY (pollId, dueTime)
);

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
    createdAt timestamp,
    
    PRIMARY KEY (pollId, answerId, voterId)
);

-- TODO: Votes by voterId