ALTER TABLE users ADD COLUMN role text DEFAULT 'user';

UPDATE users set role = 'admin' WHERE id = 1;

CREATE TABLE IF NOT EXISTS polls (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    title text NOT NULL,
    description text DEFAULT '',
    options text[] NOT NULL,
    created_by int8 NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    version integer NOT NULL DEFAULT 1
);
