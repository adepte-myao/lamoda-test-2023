CREATE TABLE storehouses (
    id TEXT PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    latitude float8 NOT NULL,
    longitude float8 NOT NULL
);

CREATE TABLE items (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    length_meters float8 NOT NULL,
    width_meters float8 NOT NULL,
    height_meters float8 NOT NULL,
    weight_kg float8 NOT NULL
);

CREATE TABLE storehouses_items (
    id BIGSERIAL PRIMARY KEY,
    storehouse_id TEXT REFERENCES storehouses (id) NOT NULL,
    item_id TEXT REFERENCES items (id) NOT NULL,
    items_count INT NOT NULL,

    CONSTRAINT storehouse_items_count_must_be_non_negative CHECK(items_count > 0),
    CONSTRAINT storehouse_and_item_ids_non_repeatable UNIQUE(storehouse_id, item_id)
);

CREATE TABLE reservations (
    id TEXT PRIMARY KEY,
    destination_latitude float8 NOT NULL,
    destination_longitude float8 NOT NULL
);

CREATE TABLE reservation_items (
    id BIGSERIAL PRIMARY KEY,
    reservation_id TEXT REFERENCES reservations (id) NOT NULL,
    item_id TEXT REFERENCES items (id) NOT NULL,
    storehouse_id TEXT REFERENCES storehouses (id) NOT NULL,
    items_count INT NOT NULL,

    CONSTRAINT reservation_items_count_must_be_non_negative CHECK(items_count > 0)
)