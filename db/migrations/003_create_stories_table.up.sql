-- Create stories table
CREATE TABLE IF NOT EXISTS stories (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    cover_image_url VARCHAR(500),
    author_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    author_name VARCHAR(100),
    status VARCHAR(20) DEFAULT 'ongoing',
    total_chapters INTEGER DEFAULT 0,
    total_views BIGINT DEFAULT 0,
    rating DECIMAL(3,2) DEFAULT 0,
    is_published BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_stories_slug ON stories(slug);
CREATE INDEX IF NOT EXISTS idx_stories_author ON stories(author_id);
CREATE INDEX IF NOT EXISTS idx_stories_status ON stories(status);

-- Create story_categories junction table (Many-to-Many)
CREATE TABLE IF NOT EXISTS story_categories (
    story_id INTEGER REFERENCES stories(id) ON DELETE CASCADE,
    category_id INTEGER REFERENCES categories(id) ON DELETE CASCADE,
    PRIMARY KEY (story_id, category_id)
);
