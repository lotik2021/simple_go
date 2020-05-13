CREATE TABLE IF NOT EXISTS maasapi.auth_regenerate_code_attempt (
    id varchar DEFAULT gen_random_uuid() PRIMARY KEY,
    device_id varchar NOT NULL,
    phone varchar NOT NULL,
    created_at timestamptz DEFAULT now(),
    nbf timestamptz NULL,
    deleted_at timestamptz NULL
);