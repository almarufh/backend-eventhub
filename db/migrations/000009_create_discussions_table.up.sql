CREATE TABLE discussions (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id INT NOT NULL,
    event_id INT,
    community_id INT NOT NULL,
    message TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_discussions_user FOREIGN KEY (user_id) 
        REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_discussions_event FOREIGN KEY (event_id) 
        REFERENCES events(id) ON DELETE CASCADE,
    CONSTRAINT fk_discussions_community FOREIGN KEY (community_id) 
        REFERENCES communities(id) ON DELETE CASCADE
);

CREATE INDEX idx_discussions_community_id ON discussions(community_id);
CREATE INDEX idx_discussions_event_id ON discussions(event_id);
CREATE INDEX idx_discussions_user_id ON discussions(user_id);
CREATE INDEX idx_discussions_created_at ON discussions(created_at);