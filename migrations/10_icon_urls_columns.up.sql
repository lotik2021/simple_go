ALTER TABLE maasapi.maas_transport_names ADD COLUMN ios_icon_url varchar;
ALTER TABLE maasapi.maas_transport_names ADD COLUMN android_icon_url varchar;

ALTER TABLE maasapi.skill DROP COLUMN IF EXISTS icon_url;

ALTER TABLE maasapi.skill ADD COLUMN ios_icon_url varchar;
ALTER TABLE maasapi.skill ADD COLUMN android_icon_url varchar;
