ALTER TABLE maasapi.user_session
    ADD COLUMN origin_place_id VARCHAR,
    ADD COLUMN destination_place_id VARCHAR;