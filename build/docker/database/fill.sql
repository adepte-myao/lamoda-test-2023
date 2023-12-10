INSERT INTO storehouses (id, name, latitude, longitude)
VALUES ('a', 'a', 50, 50),
       ('b', 'b', 55, 60),
       ('c', 'c', 60, 40),
       ('d', 'd', 65, 50),
       ('e', 'e', 70, 60),
       ('f', 'f', 45, 40),
       ('g', 'g', 40, 50),
       ('h', 'h', 35, 60),
       ('i', 'i', 30, 40),
       ('j', 'j', 25, 50);

INSERT INTO items (id, name, length_meters, width_meters, height_meters, weight_kg)
SELECT  series.series::text, series.series::text, random()*20, random()*10, random()*5, random()*50 FROM generate_series(1, 20) AS series;

INSERT INTO storehouses_items (storehouse_id, item_id, items_count)
SELECT storehouses.id, items.id, trunc(random()*9+7) FROM storehouses JOIN items ON true WHERE random() < 0.8;

INSERT INTO reservations (id, destination_latitude, destination_longitude)
VALUES ('one-reservation', 42, 43), ('two-reservation', 48, 49);

INSERT INTO reservation_items (reservation_id, item_id, storehouse_id, items_count)
SELECT 'one-reservation', si.item_id, si.storehouse_id, si.items_count-3
FROM storehouses_items AS si ORDER BY random() LIMIT 3;

INSERT INTO reservation_items (reservation_id, item_id, storehouse_id, items_count)
SELECT 'two-reservation', si.item_id, si.storehouse_id, si.items_count-3
FROM storehouses_items AS si ORDER BY random() LIMIT 3;