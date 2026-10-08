
CREATE TABLE saved_events_users (
    user_id INT NOT NULL,
    event_id INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, event_id),
    CONSTRAINT fk_saved_events_users_user FOREIGN KEY (user_id) 
        REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_saved_events_users_event FOREIGN KEY (event_id) 
        REFERENCES events(id) ON DELETE CASCADE
);

CREATE INDEX idx_joined_events_users_event_id ON joined_events_users(event_id);
CREATE INDEX idx_saved_events_users_event_id ON saved_events_users(event_id);