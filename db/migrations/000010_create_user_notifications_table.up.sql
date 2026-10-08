CREATE TABLE user_notifications (
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    notification_id INT NOT NULL REFERENCES notifications(id) ON DELETE CASCADE,
    is_read BOOLEAN DEFAULT FALSE,
    read_at TIMESTAMPTZ,
    PRIMARY KEY (user_id, notification_id)
);

CREATE INDEX idx_user_notifications_notification_id ON user_notifications(notification_id);