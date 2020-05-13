CREATE TABLE IF NOT EXISTS maasapi.user_search_result (
    id varchar not null unique primary key,
    version varchar not null,
    user_id varchar references maasapi.user_profile("id"),
    distance_in_km int not null,
    origin geography not null,
	destination geography not null,
	origin_google_place_id varchar references maasapi.google_place ("google_place_id") on update cascade,
	destination_google_place_id varchar references maasapi.google_place ("google_place_id") on update cascade,
	origin_movista_place_id varchar,
	destination_movista_place_id varchar,
	origin_movista_place_name varchar,
	destination_movista_place_name varchar,
    departure_time timestamptz,
    arrival_time timestamptz,
    response jsonb not null,
    is_error boolean,
    error varchar,
    created_at timestamptz not null default now()
);

CREATE INDEX idx_user_search_result_response ON maasapi.user_search_result USING GIN (response);
CREATE INDEX idx_user_search_result_response_ids ON maasapi.user_search_result USING GIN ((response -> 'ids'));
CREATE INDEX idx_user_search_result_origin ON maasapi.user_search_result USING GIST (origin);
CREATE INDEX idx_user_search_result_destination ON maasapi.user_search_result USING GIST (destination);