CREATE TABLE message {
    id         INT,
    content    TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
}