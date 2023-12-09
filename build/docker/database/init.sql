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
    weight_kg float8 NOT NULL,
    count INT NOT NULL,

    CONSTRAINT items_count_must_be_non_negative CHECK(count > 0)
)