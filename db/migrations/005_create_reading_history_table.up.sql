-- Create reading_history table
CREATE TABLE IF NOT EXISTS reading_history (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    story_id INTEGER REFERENCES stories(id) ON DELETE CASCADE,
    last_chapter_id INTEGER REFERENCES chapters(id) ON DELETE SET NULL,
    last_read_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, story_id)
);

CREATE INDEX IF NOT EXISTS idx_reading_history_user ON reading_history(user_id);
