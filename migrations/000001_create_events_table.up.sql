CREATE TABLE events (
    ID SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NULL,
    location TEXT NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    price DECIMAL NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Constraints

-- Indexes
CREATE INDEX events_title_index ON events (title);
CREATE INDEX events_start_time_index ON events (start_time);