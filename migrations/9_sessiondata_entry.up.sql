ALTER TABLE maasapi.session_data ADD COLUMN user_entry_data jsonb default null;
ALTER TABLE maasapi.session_data ADD COLUMN actions jsonb default null;
ALTER TABLE maasapi.session_data ADD COLUMN objects jsonb default null;

ALTER TABLE maasapi.user_session ADD COLUMN trip_type varchar default null;