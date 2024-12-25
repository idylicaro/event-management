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

