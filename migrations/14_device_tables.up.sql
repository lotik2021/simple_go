ALTER TABLE IF EXISTS maasapi.google_place RENAME COLUMN google_place_id TO id;
ALTER TABLE IF EXISTS maasapi.google_place ADD COLUMN IF NOT EXISTS created_at timestamptz DEFAULT now();
ALTER TABLE IF EXISTS maasapi.google_place ADD COLUMN IF NOT EXISTS updated_at timestamptz;

ALTER TABLE IF EXISTS maasapi.skill ALTER COLUMN id TYPE varchar USING gen_random_uuid();
ALTER TABLE IF EXISTS maasapi.skill ALTER COLUMN id SET default gen_random_uuid();
ALTER TABLE IF EXISTS maasapi.skill DROP CONSTRAINT IF EXISTS skills_title_key;
ALTER TABLE IF EXISTS maasapi.skill ADD CONSTRAINT skills_title_disable_in_versions_key UNIQUE (title, disable_in_versions);
ALTER TABLE IF EXISTS maasapi.skill ADD COLUMN IF NOT EXISTS created_at timestamptz DEFAULT now();
DROP SEQUENCE IF EXISTS maasapi.skills_id_seq;

CREATE TABLE IF NOT EXISTS maasapi.device (
    id varchar DEFAULT gen_random_uuid() PRIMARY KEY,
    "name" varchar NULL,
    last_name varchar NULL,
    device_os varchar NULL,
    device_category varchar NULL,
    device_type varchar NULL,
    device_info varchar NULL,
    os_player_id varchar NULL,
    name_asked_at timestamptz NULL,
    google_transit_modes varchar[] NULL,
    created_at timestamptz DEFAULT now(),
    updated_at timestamptz DEFAULT now(),
    deleted_at timestamptz NULL
);

INSERT INTO maasapi.device(id, "name", last_name, os_player_id, name_asked_at, created_at)
    SELECT id, user_name, last_name, os_player_id, last_name_question_time, created_at from maasapi.user_profile;

UPDATE maasapi.device d set google_transit_modes = (
    select array_agg(google_transport) from maasapi.user_google_transport where user_id = d.id and filter_status = true
);

UPDATE maasapi.device SET google_transit_modes = '{bus,train,tram,subway}' where google_transit_modes is NULL;


ALTER TABLE IF EXISTS maasapi.user_location DROP CONSTRAINT user_location_user_id_fkey;
ALTER TABLE IF EXISTS maasapi.user_location DROP CONSTRAINT user_location_pkey;
ALTER TABLE IF EXISTS maasapi.user_location DROP COLUMN id;
ALTER TABLE IF EXISTS maasapi.user_location RENAME TO device_location;
ALTER TABLE IF EXISTS maasapi.device_location RENAME COLUMN user_id TO device_id;
CREATE INDEX IF NOT EXISTS device_location_pkey ON maasapi.device_location USING btree (device_id, created_at);
DROP SEQUENCE IF EXISTS maasapi.user_location_id_seq;


ALTER TABLE IF EXISTS maasapi.session_data DROP CONSTRAINT session_data_session_id_fkey;
ALTER TABLE IF EXISTS maasapi.session_data RENAME TO device_session_data;
ALTER TABLE IF EXISTS maasapi.device_session_data ALTER COLUMN id TYPE varchar;


ALTER TABLE IF EXISTS maasapi.user_session DROP CONSTRAINT user_session_user_id_fkey;
ALTER TABLE IF EXISTS maasapi.user_session RENAME TO device_session;
ALTER TABLE IF EXISTS maasapi.device_session ALTER COLUMN id SET default gen_random_uuid();
ALTER TABLE IF EXISTS maasapi.device_session ADD COLUMN IF NOT EXISTS updated_at timestamptz;
ALTER TABLE IF EXISTS maasapi.device_session RENAME COLUMN user_id TO device_id;


ALTER TABLE IF EXISTS maasapi.taxi_order DROP CONSTRAINT taxi_order_user_id_fkey;
ALTER TABLE IF EXISTS maasapi.taxi_order RENAME TO device_taxi_order;
ALTER TABLE IF EXISTS maasapi.device_taxi_order RENAME COLUMN user_id TO device_id;
ALTER TABLE IF EXISTS maasapi.device_taxi_order DROP COLUMN device_os;
ALTER TABLE IF EXISTS maasapi.device_taxi_order ADD COLUMN IF NOT EXISTS updated_at timestamptz;


ALTER TABLE IF EXISTS maasapi.travel_card_payment DROP CONSTRAINT travel_card_payment_user_id_fkey;
ALTER TABLE IF EXISTS maasapi.travel_card_payment RENAME TO device_travel_card_payment;
ALTER TABLE IF EXISTS maasapi.device_travel_card_payment RENAME COLUMN user_id TO device_id;
ALTER TABLE IF EXISTS maasapi.device_travel_card_payment ALTER COLUMN id TYPE varchar;
ALTER TABLE IF EXISTS maasapi.device_travel_card_payment ALTER COLUMN currency TYPE varchar;
ALTER TABLE IF EXISTS maasapi.device_travel_card_payment ALTER COLUMN currency SET default 'rub';
ALTER TABLE IF EXISTS maasapi.device_travel_card_payment ADD COLUMN IF NOT EXISTS updated_at timestamptz;


ALTER TABLE IF EXISTS maasapi.user_favorite DROP CONSTRAINT user_favorite_place_id_fkey;
ALTER TABLE IF EXISTS maasapi.user_favorite DROP CONSTRAINT user_favorite_user_id_fkey;
ALTER TABLE IF EXISTS maasapi.user_favorite RENAME TO device_favorite;
ALTER TABLE IF EXISTS maasapi.device_favorite RENAME COLUMN user_id TO device_id;
ALTER TABLE IF EXISTS maasapi.device_favorite ALTER COLUMN "type" TYPE varchar;
ALTER TABLE IF EXISTS maasapi.device_favorite ADD COLUMN IF NOT EXISTS created_at timestamptz DEFAULT now();
ALTER TABLE IF EXISTS maasapi.device_favorite ADD COLUMN IF NOT EXISTS updated_at timestamptz DEFAULT now();
ALTER SEQUENCE maasapi.user_favorite_id_seq RENAME TO device_favorite_id_seq;
ALTER TABLE IF EXISTS maasapi.device_favorite RENAME CONSTRAINT user_favorite_pkey TO device_favorite_pkey;

ALTER TABLE IF EXISTS maasapi.user_history DROP CONSTRAINT user_history_place_id_fkey;
ALTER TABLE IF EXISTS maasapi.user_history DROP CONSTRAINT user_history_user_id_fkey;
ALTER TABLE IF EXISTS maasapi.user_history RENAME TO device_google_place_history;
ALTER TABLE IF EXISTS maasapi.device_google_place_history ALTER COLUMN id TYPE varchar USING gen_random_uuid();
ALTER TABLE IF EXISTS maasapi.device_google_place_history ALTER COLUMN id SET default gen_random_uuid();
DROP SEQUENCE IF EXISTS maasapi.user_history_id_seq;
ALTER TABLE IF EXISTS maasapi.device_google_place_history RENAME COLUMN user_id TO device_id;
ALTER TABLE IF EXISTS maasapi.device_google_place_history ADD COLUMN IF NOT EXISTS created_at timestamptz DEFAULT now();
ALTER TABLE IF EXISTS maasapi.device_google_place_history ADD CONSTRAINT device_google_place_history_device_id_place_id UNIQUE (device_id, place_id);

CREATE TABLE IF NOT EXISTS maasapi.device_movista_search_history (
    id varchar DEFAULT gen_random_uuid() PRIMARY KEY,
    device_id varchar NOT NULL,
    origin geography NULL,
    destination geography NULL,
	from_id integer NOT NULL,
	to_id integer NOT NULL,
	from_google_place_id varchar NULL,
	to_google_place_id varchar NULL,
	trip_types varchar[] NULL,
	departure_time timestamptz NOT NULL,
	arrival_time timestamptz NULL,
    created_at timestamptz DEFAULT now()
);

INSERT INTO maasapi.device_movista_search_history (id, device_id, origin, destination, from_id, to_id, from_google_place_id, to_google_place_id, departure_time, arrival_time, created_at)
    SELECT id, user_id, origin, destination, origin_movista_place_id::integer, destination_movista_place_id::integer, origin_google_place_id, destination_google_place_id, departure_time, arrival_time, created_at
    from maasapi.user_search_result where origin_movista_place_id is not NULL;

INSERT INTO maasapi.device_movista_search_history (id, device_id, from_id, to_id, trip_types, departure_time, created_at)
    SELECT id, user_id, from_id, to_id, trip_types, departure_time, created_at
    from maasapi.movista_user_search_history;

CREATE TABLE IF NOT EXISTS maasapi.device_google_search_history (
    id varchar DEFAULT gen_random_uuid() PRIMARY KEY,
    device_id varchar NOT NULL,
    origin geography NOT NULL,
    destination geography NOT NULL,
    origin_place_id varchar NULL,
    destination_place_id varchar NULL,
    departure_time timestamptz NULL,
    arrival_time timestamptz NULL,
    trip_types varchar[] NULL,
    created_at timestamptz DEFAULT now()
);

INSERT INTO maasapi.device_google_search_history (id, device_id, origin, destination, origin_place_id, destination_place_id, departure_time, arrival_time, created_at)
    SELECT id, user_id, origin, destination, origin_google_place_id, destination_google_place_id, departure_time, arrival_time, created_at
    from maasapi.user_search_result where origin_movista_place_id is NULL;


ALTER TABLE IF EXISTS maasapi.user_google_transport DROP CONSTRAINT user_google_transport_google_transport_fkey;
ALTER TABLE IF EXISTS maasapi.maas_transport_names RENAME TO google_transit_mode_name;
ALTER TABLE IF EXISTS maasapi.google_transit_mode_name RENAME COLUMN google_transport TO "name";
ALTER TABLE IF EXISTS maasapi.google_transit_mode_name RENAME COLUMN transport_name TO display_name;
ALTER TABLE IF EXISTS maasapi.google_transit_mode_name DROP CONSTRAINT maas_transport_names_pkey;
ALTER TABLE IF EXISTS maasapi.google_transit_mode_name ALTER COLUMN "name" TYPE varchar;
ALTER TABLE IF EXISTS maasapi.google_transit_mode_name ADD CONSTRAINT google_transit_mode_name_pkey PRIMARY KEY ("name");