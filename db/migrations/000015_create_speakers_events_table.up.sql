CREATE TABLE speakers_events (
    speaker_id INT NOT NULL,
    event_id INT NOT NULL,
    PRIMARY KEY (speaker_id, event_id),
    CONSTRAINT fk_speakers_events_speaker FOREIGN KEY (speaker_id) 
        REFERENCES speakers(id) ON DELETE CASCADE,
    CONSTRAINT fk_speakers_events_event FOREIGN KEY (event_id) 
        REFERENCES events(id) ON DELETE CASCADE
);

CREATE INDEX idx_speakers_events_event_id ON speakers_events(event_id);