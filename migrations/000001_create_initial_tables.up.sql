CREATE TYPE pull_request_status AS ENUM ('OPEN', 'MERGED');

CREATE TABLE IF NOT EXISTS teams
(
    name VARCHAR(255) PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS users
(
    id        VARCHAR(255) PRIMARY KEY,
    username  VARCHAR(255) NOT NULL,
    team_name VARCHAR(255) NOT NULL REFERENCES teams (name),
    is_active BOOLEAN      NOT NULL
);

CREATE TABLE IF NOT EXISTS pull_requests
(
    id         VARCHAR(255) PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    author_id  VARCHAR(255) NOT NULL REFERENCES users (id),
    status     pull_request_status DEFAULT 'OPEN',
    created_at TIMESTAMPTZ  NOT NULL,
    merged_at  TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS review_assignments
(
    user_id         VARCHAR(255) REFERENCES users (id),
    pull_request_id VARCHAR(255) REFERENCES pull_requests (id),

    PRIMARY KEY (user_id, pull_request_id)
);

CREATE INDEX idx_users_team_active ON users (team_name, is_active);
CREATE INDEX idx_pull_requests_author ON pull_requests (author_id);
