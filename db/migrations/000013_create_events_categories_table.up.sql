CREATE TABLE events_categories (
    category_id INT NOT NULL,
    event_id INT NOT NULL,
    PRIMARY KEY (category_id, event_id),
    CONSTRAINT fk_events_categories_category FOREIGN KEY (category_id) 
        REFERENCES categories(id) ON DELETE CASCADE,
    CONSTRAINT fk_events_categories_event FOREIGN KEY (event_id) 
        REFERENCES events(id) ON DELETE CASCADE
);

CREATE INDEX idx_events_categories_event_id ON events_categories(event_id);