CREATE TABLE IF NOT EXISTS maasapi.movista_user_search_history (
    id varchar,
    user_id varchar not null,
	from_id integer not null,
	to_id integer not null,
	trip_types varchar[] default null,
	departure_time timestamptz not null,
    created_at timestamptz not null default now()
);