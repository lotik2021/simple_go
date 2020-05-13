CREATE SCHEMA IF NOT EXISTS maasapi;
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TYPE google_transport_enum AS ENUM ('bus', 'subway', 'tram', 'train');
CREATE TYPE favorite_place_type AS ENUM ('home', 'work', 'none');
CREATE TYPE currency AS ENUM ('rub');

CREATE TABLE IF NOT EXISTS maasapi.user_profile (
    id text DEFAULT gen_random_uuid(),
    user_name text,
    created_at timestamptz not null default now(),
    auth_user_id varchar,
    auth_device_id varchar,
    last_name varchar,
    last_name_question_time timestamptz,
    primary key (id)
);

CREATE TABLE IF NOT EXISTS maasapi.google_place (
    google_place_id varchar not null unique,
    main_text varchar,
    secondary_text varchar,
    coordinate geography NOT NULL,
    primary key (google_place_id)
);

CREATE TABLE IF NOT EXISTS maasapi.maas_transport_names
(
    google_transport google_transport_enum not null primary key,
    transport_name varchar(30) not null,
    icon_name varchar not null
);

CREATE TABLE IF NOT EXISTS maasapi.user_google_transport
(
    user_id varchar references maasapi.user_profile("id") not null,
    google_transport google_transport_enum references maasapi.maas_transport_names("google_transport") not null,
    filter_status boolean not null,
    primary key (user_id, google_transport)
);

CREATE TABLE IF NOT EXISTS maasapi.metro (
    id SERIAL NOT null PRIMARY KEY,
	"language" varchar(5) NULL,
	city_name varchar(100) NULL,
	movista_city_id int4 NULL,
	logo_uri varchar(255) NULL,
	line_number varchar(25) NULL,
	full_name varchar(80) NULL,
	short_name varchar(40) NULL,
	color_hex varchar(50) NULL,
	agency_url varchar(80) NULL,
	agency_name varchar(100) NULL,
	"type" varchar(15) NULL
);

CREATE TABLE IF NOT EXISTS maasapi.user_favorite (
     id SERIAL NOT null PRIMARY KEY,
     user_id varchar references maasapi.user_profile("id") NOT NULL,
     place_id varchar references maasapi.google_place("google_place_id") NOT NULL,
     type favorite_place_type NOT NULL,
     updated_at timestamptz not null default now(),
     name varchar
);

CREATE TABLE IF NOT EXISTS maasapi.user_history (
     id SERIAL NOT null PRIMARY KEY,
     user_id varchar references maasapi.user_profile("id") NOT NULL,
     place_id varchar references maasapi.google_place("google_place_id") NOT NULL,
     number_of_searches int default 0,
     updated_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS maasapi.user_location (
     id SERIAL not null primary key,
     user_id varchar references maasapi.user_profile("id") NOT NULL,
     coordinate geography NOT NULL,
     created_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS maasapi.taxi_order (
	id uuid DEFAULT gen_random_uuid() not null primary key,
	user_id varchar references maasapi.user_profile("id") NOT null,
	provider varchar NOT NULL,
	device_os varchar,
	used_link varchar,
	created_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS maasapi.travel_card_payment (
	id uuid DEFAULT gen_random_uuid() not null primary key,
	user_id varchar references maasapi.user_profile("id") NOT null,
	amount float not null,
	card_number varchar not null default '',
	card_type varchar not null,
	processed bool default false,
	currency currency default 'rub' not null,
	created_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS maasapi.skill (
    id SERIAL not null primary key,
	skill_description varchar not null unique,
	ord int not null unique,
	icon varchar not null,
	action_id varchar not null
);

CREATE TABLE IF NOT EXISTS maasapi.application_version (
	id varchar not null primary key
);

CREATE TABLE IF NOT EXISTS maasapi.skill_version (
	skill_id SERIAL references maasapi.skill("id") not null,
	version_id varchar references maasapi.application_version("id") NOT null
);

CREATE TABLE IF NOT EXISTS maasapi.user_session (
    id varchar not null primary key,
	user_id varchar references maasapi.user_profile("id") not null,
	origin geography,
	destination geography,
    departure_time timestamptz,
    arrival_time timestamptz,
	state varchar not null,
	created_at timestamptz not null default now()
);

CREATE TABLE IF NOT EXISTS maasapi.session_data (
	id uuid not null DEFAULT gen_random_uuid() primary key,
	session_id varchar references maasapi.user_session("id") not null,
	created_at timestamptz not null default now(),
	action_id varchar default null,
    action_name varchar not null,
    user_response varchar default null,
    dialog_response varchar default null
);

INSERT INTO maasapi.application_version(id) VALUES
('1.1')
,('1.2')
;

INSERT INTO maasapi.skill(id, skill_description, ord, icon, action_id) VALUES
(1, E'Найти билеты на\nпоезд, самолет, автобус', 1, 'icon_world_route', 'build_route')
,(2, E'Проложить маршруты\nпо городу и области',  2, 'icon_local_route', 'build_route')
,(3, E'Заказать такси',  3, 'icon_taxi_route', 'build_route')
,(4, E'Пополнить транспортные\nкарты: тройку и стрелку',  4, 'icon_travel_card', 'card')
;

INSERT INTO maasapi.skill_version(skill_id, version_id) VALUES
(1, '1.1')
,(2, '1.1')
,(3, '1.1')
,(4, '1.1')
,(4, '1.2')
;

INSERT INTO maasapi.maas_transport_names (google_transport,transport_name,icon_name) VALUES
('bus','Автобус','icon_bus')
,('subway','Метро','icon_metro')
,('tram','Трамвай','icon_tram')
,('train','Электричка','icon_commuter_train');

INSERT INTO maasapi.metro ("language",city_name,movista_city_id,logo_uri,line_number,full_name,short_name,color_hex,agency_url,agency_name,"type") VALUES
('ru','Москва',65537,'','1','Сокольническая линия','Сокольническая','#E01E21','http://mosmetro.ru/','ГУП "Московский метрополитен"','subway')
,('ru','Москва',65537,'','2','Замоскворецкая линия','Замоскворецкая','#46BB49','http://mosmetro.ru/','ГУП "Московский метрополитен"','subway')
,('ru','Москва',65537,'','3','Арбатско-Покровская линия','Арбатско-Покровская','#007BCF','http://mosmetro.ru/','ГУП "Московский метрополитен"','subway')
,('ru','Москва',65537,'','4','Филёвская линия','Филёвская','#02BAFF','http://mosmetro.ru/','ГУП "Московский метрополитен"','subway')
,('ru','Москва',65537,'','5','Кольцевая линия','Кольцевая','#905035','http://mosmetro.ru/','ГУП "Московский метрополитен"','subway')
,('ru','Москва',65537,'','6','Калужско-Рижская линия','Калужско-Рижская','#FF8335','http://mosmetro.ru/','ГУП "Московский метрополитен"','subway')
,('ru','Москва',65537,'','7','Таганско-Краснопресненская линия','Таганско-Краснопресненская','#9D389C','http://mosmetro.ru/','ГУП "Московский метрополитен"','subway')
,('ru','Москва',65537,'','8','Калининская линия','Калининская ','#FFD600','http://mosmetro.ru/','ГУП "Московский метрополитен"','subway')
,('ru','Москва',65537,'','8A','Солнцевская линия','Солнцевская','#FFD600','http://mosmetro.ru/','ГУП "Московский метрополитен"','subway')
,('ru','Москва',65537,'','9','Серпуховско-Тимирязевская линия','Серпуховско-Тимирязевская','#ABA9AA','http://mosmetro.ru/','ГУП "Московский метрополитен"','subway')
,('ru','Москва',65537,'','10','Люблинско-Дмитровская линия','Люблинско-Дмитровская','#C7E31D','http://mosmetro.ru/','ГУП "Московский метрополитен"','subway')
,('ru','Москва',65537,'','11','Большая Кольцевая линия','Большая Кольцевая','#81D6DA','http://mosmetro.ru/','ГУП "Московский метрополитен"','subway')
,('ru','Москва',65537,'','11A','Каховская линия','Каховская','#81D6DA','http://mosmetro.ru/','ГУП "Московский метрополитен"','subway')
,('ru','Москва',65537,'','12','Бутовская линия','Бутовская','#BBC6E8','http://mosmetro.ru/','ГУП "Московский метрополитен"','subway')
,('ru','Москва',65537,'','13','Московский монорельс','Московский монорельс','#BBC6E9','http://mosmetro.ru/','ГУП "Московский метрополитен"','subway')
,('ru','Москва',65537,'','14','Московское Центральное Кольцо','МЦК','#FFB7BA','http://mosmetro.ru/','ГУП "Московский метрополитен"','subway')
,('ru','Москва',65537,'','15','Некрасовская линия','Некрасовская','#ED89B7','http://mosmetro.ru/','ГУП "Московский метрополитен"','subway')
,('ru','Санкт-Петербург',83130,'','1','Кировско-Выборгская линия','Кировско-Выборгская','#E01E21','http://www.metro.spb.ru','ГУП "Петербургский метрополитен"','subway')
,('ru','Санкт-Петербург',83130,'','2','Московско-Петроградская линия','Московско-Петроградская','#007BCF','http://www.metro.spb.ru','ГУП "Петербургский метрополитен"','subway')
,('ru','Санкт-Петербург',83130,'','3','Невско-Василеостровская линия','Невско-Василеостровская','#46BB49','http://www.metro.spb.ru','ГУП "Петербургский метрополитен"','subway')
,('ru','Санкт-Петербург',83130,'','4','Правобережная линия','Правобережная','#FF8335','http://www.metro.spb.ru','ГУП "Петербургский метрополитен"','subway')
,('ru','Санкт-Петербург',83130,'','5','Фрунзенско-Приморская линия','Фрунзенско-Приморская','#9D389C','http://www.metro.spb.ru','ГУП "Петербургский метрополитен"','subway')
,('ru','Нижний-Новгород',90702,'','Сормовская','Сормовско-Мещерская линия','Сормовско-Мещерская','#007BCF','http://metronn.ru/','МУП "Нижегородское метро"','subway')
,('ru','Нижний-Новгород',90702,'','Автозаводская','Автозаводская линия','Автозаводская','#E01E21','http://metronn.ru/','МУП "Нижегородское метро"','subway')
,('ru','Новосибирск',74177,'','Ленинская','Ленинская линия','Ленинская','#E01E21','http://www.nsk-metro.ru/','Новосибирский метрополитен','subway')
,('ru','Новосибирск',74177,'','Дзержинская','Дзержинская линия','Дзержинская','#46BB49','http://www.nsk-metro.ru/','Новосибирский метрополитен','subway')
,('ru','Самара',86171,'','Первая линия','Первая линия','Первая линия','#E01E21','http://metrosamara.ru/','МП г.о. Самара "Самарский метрополитен"','subway')
,('ru','Волгоград',79603,'','СТ','СТ','СТ','#E01E21','http://www.gortransvolga.ru/','Метроэлектротранс','TRAM')
,('ru','Волгоград',79603,'','СТ2','СТ','СТ','#E01E22','http://www.gortransvolga.ru/','Метроэлектротранс','TRAM')
,('ru','Екатеринбург',77762,'','Первая линия','Первая линия','Первая линия','#46BB49','http://metro-ektb.ru/','ЕМУП "Екатеринбургский Метрополитен"','subway')
,('ru','Казань',83781,'','Центральная','Центральная линия','Центральная','#E01E21','http://kazanmetro.ru/','МУП "Метроэлектротранс"','subway')
;



