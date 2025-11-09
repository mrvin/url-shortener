CREATE TABLE IF NOT EXISTS users (
	name TEXT NOT NULL UNIQUE PRIMARY KEY,
	hash_password TEXT NOT NULL,
	role TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS url(
	alias TEXT NOT NULL UNIQUE PRIMARY KEY,
	url TEXT NOT NULL,
	count BIGINT NOT NULL CHECK (count >= 0) DEFAULT 0,
	user_name TEXT references users(name) on delete cascade,
	created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_url_user_name ON url(user_name);

CREATE OR REPLACE FUNCTION get_url(alias_for_url TEXT)
RETURNS TEXT
LANGUAGE plpgsql
AS $$
DECLARE
    full_url TEXT;
BEGIN
    UPDATE url SET count = count+1 WHERE alias = alias_for_url RETURNING url INTO full_url;
    IF NOT FOUND THEN
        RAISE EXCEPTION 'alias % not found', alias_for_url;
    END IF;
    RETURN full_url;
END;
$$;
