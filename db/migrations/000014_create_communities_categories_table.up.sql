CREATE TABLE communities_categories (
    category_id INT NOT NULL,
    community_id INT NOT NULL,
    PRIMARY KEY (category_id, community_id),
    CONSTRAINT fk_communities_categories_category FOREIGN KEY (category_id) 
        REFERENCES categories(id) ON DELETE CASCADE,
    CONSTRAINT fk_communities_categories_community FOREIGN KEY (community_id) 
        REFERENCES communities(id) ON DELETE CASCADE
);

CREATE INDEX idx_communities_categories_community_id ON communities_categories(community_id);