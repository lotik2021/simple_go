CREATE TABLE IF NOT EXISTS maasapi.user_device_session (
    id SERIAL PRIMARY KEY,
    device_id varchar NOT NULL,
    user_id integer NOT NULL,
    created_at timestamptz DEFAULT now(),
    deleted_at timestamptz NULL
);

INSERT INTO maasapi.user_device_session(device_id, user_id)
    SELECT id, auth_user_id::int from maasapi.user_profile where auth_user_id is not NULL and auth_user_id != '';