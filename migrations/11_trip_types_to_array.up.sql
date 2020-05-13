ALTER TABLE maasapi.user_session DROP COLUMN trip_type;

ALTER TABLE maasapi.user_session ADD COLUMN trip_types varchar[] default null;