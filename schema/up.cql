CREATE TABLE IF NOT EXISTS polls (
    poll_id uuid,
    title text,
    description text,
    poll_type int,
    due_time timestamp,

    PRIMARY KEY (poll_id, due_time)
) WITH CLUSTERING ORDER BY (due_time DESC);

CREATE TABLE IF NOT EXISTS active_polls (
    poll_id uuid,
    title text,
    description text,
    poll_type int,
    due_time timestamp,

    PRIMARY KEY (poll_id, due_time)
) WITH CLUSTERING ORDER BY (due_time DESC);

CREATE TABLE IF NOT EXISTS answers (
    answer_id uuid,
    answer text,
    poll_id uuid,

    PRIMARY KEY (poll_id, answer_id)
);

CREATE TABLE IF NOT EXISTS votes (
    poll_id uuid,
    answer_id uuid,
    voter_id uuid,
    
    PRIMARY KEY (poll_id, answer_id, voter_id)
);
