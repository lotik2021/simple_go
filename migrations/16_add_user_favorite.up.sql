CREATE TABLE IF NOT EXISTS maasapi.user_favorite (
    id SERIAL PRIMARY KEY,
    user_id integer NOT NULL,
    device_id varchar NOT NULL,
    place_id varchar NOT NULL,
    type varchar NOT NULL,
    name varchar NULL,
    in_actions boolean DEFAULT false,
    created_at timestamptz DEFAULT now(),
    updated_at timestamptz DEFAULT now()
);