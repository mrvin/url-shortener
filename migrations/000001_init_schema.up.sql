CREATE TABLE IF NOT EXISTS users (
	name TEXT NOT NULL UNIQUE PRIMARY KEY,
	hash_password TEXT NOT NULL,
	role TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS urls(
	alias TEXT NOT NULL UNIQUE PRIMARY KEY,
	url TEXT NOT NULL,
	count BIGINT NOT NULL CHECK (count >= 0) DEFAULT 0,
	username TEXT references users(name) on delete cascade,
	created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_urls_username ON urls(username);
