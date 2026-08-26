CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY
);

INSERT INTO users (id)
VALUES
    ('user-1'),
    ('user-2'),
    ('user-3')
ON CONFLICT (id) DO NOTHING;