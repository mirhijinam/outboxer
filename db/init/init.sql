CREATE TABLE message (
    id         SERIAL PRIMARY KEY,
    content    TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE event (
    id           SERIAL PRIMARY KEY,
    payload      TEXT NOT NULL,
    status       TEXT NOT NULL DEFAULT 'new' CHECK(status IN ('new', 'done')),
    created_at   TIMESTAMP NOT NULL DEFAULT NOW(),
    reserved_for TIMESTAMP DEFAULT NOT NULL
);