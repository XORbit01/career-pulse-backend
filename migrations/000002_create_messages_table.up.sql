CREATE TABLE IF NOT EXISTS conversations (
    id SERIAL PRIMARY KEY,
    participant_one_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    participant_two_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (participant_one_id, participant_two_id)
);

CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    conversation_id INT NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    sender_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
