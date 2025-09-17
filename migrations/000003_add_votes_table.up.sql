CREATE TABLE IF NOT EXISTS votes (
    id bigserial PRIMARY KEY,
    poll_id int8 NOT NULL REFERENCES polls(id) ON DELETE CASCADE,
    user_id int8 NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    chosen_option text NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),

    CONSTRAINT user_poll_vote_unique UNIQUE (poll_id, user_id)
);
