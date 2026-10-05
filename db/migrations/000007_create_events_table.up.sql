CREATE TABLE events (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    community_id INT NOT NULL,
    organizer_id INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    location VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    image VARCHAR(255) NOT NULL,
    capacity INT NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_events_community FOREIGN KEY (community_id) 
        REFERENCES communities(id) ON DELETE CASCADE,
    CONSTRAINT fk_events_organizer FOREIGN KEY (organizer_id) 
        REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_events_community_id ON events(community_id);
CREATE INDEX idx_events_organizer_id ON events(organizer_id);
CREATE INDEX idx_events_start_time ON events(start_time);