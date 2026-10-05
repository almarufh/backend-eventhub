CREATE TABLE profiles (
    user_id INT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'attendee' 
        CHECK (role IN ('organizer', 'attendee', 'admin')),
    dark_preference BOOLEAN DEFAULT FALSE,
    address TEXT,
    job VARCHAR(255),
    office VARCHAR(255),
    image VARCHAR(255),
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_profiles_user FOREIGN KEY (user_id) 
        REFERENCES users(id) ON DELETE CASCADE
);