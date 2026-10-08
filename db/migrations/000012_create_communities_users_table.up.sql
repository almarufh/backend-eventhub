CREATE TABLE communities_users (
    user_id INT NOT NULL,
    community_id INT NOT NULL,
    PRIMARY KEY (user_id, community_id),
    CONSTRAINT fk_communities_users_user FOREIGN KEY (user_id) 
        REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_communities_users_community FOREIGN KEY (community_id) 
        REFERENCES communities(id) ON DELETE CASCADE
);

CREATE INDEX idx_communities_users_community_id ON communities_users(community_id);