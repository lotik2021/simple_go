CREATE TABLE IF NOT EXISTS maasapi.citizenship
(
    id integer NOT NULL,
    smallname character varying(100) NOT NULL,
    fullname character varying(200),
    alpha2 character varying(2) NOT NULL,
    alpha3 character varying(3) NOT NULL,
    lng character varying(2) NOT NULL,
    CONSTRAINT citizenship_pkey PRIMARY KEY (id)
);
