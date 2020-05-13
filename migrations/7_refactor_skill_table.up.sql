DROP TABLE IF EXISTS maasapi.skill_version;
DROP TABLE IF EXISTS maasapi.application_version;

CREATE TABLE IF NOT EXISTS maasapi.skills (
    id SERIAL not null primary key,
	title varchar not null unique,
	ord int not null unique,
	icon varchar,
	icon_url varchar,
	action_id varchar not null,
	disable_in_versions varchar
);

INSERT INTO maasapi.skills (id, title, ord, icon, action_id)
SELECT id, skill_description, ord, icon, action_id from maasapi.skill;

DROP TABLE IF EXISTS maasapi.skill;
ALTER TABLE maasapi.skills RENAME TO skill;