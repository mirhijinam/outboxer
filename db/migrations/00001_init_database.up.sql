CREATE TABLE IF NOT EXISTS "message" (
    id         INT GENERATED ALWAYS AS IDENTITY NOT NULL,
    content    TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS "event" (
    id             INT GENERATED ALWAYS AS IDENTITY NOT NULL,
    payload        TEXT NOT NULL,
    status         TEXT NOT NULL DEFAULT 'new' CHECK(status IN ('new', 'done')),
    created_at     TIMESTAMP NOT NULL DEFAULT NOW(),
    reserved_until TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);
