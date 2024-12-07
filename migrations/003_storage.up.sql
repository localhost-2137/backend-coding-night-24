CREATE TABLE items
(
    id             INTEGER PRIMARY KEY,
    name           TEXT NOT NULL,
    description    TEXT,
    quantity       INTEGER   DEFAULT 1,
    unit           TEXT,
    critical_level INTEGER   DEFAULT 0,
    last_checked   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    location       TEXT
);

INSERT INTO items (name, description, quantity, unit, critical_level, location)
VALUES ('Apteczka', 'Zestaw pierwszej pomocy', 10, 'szt.', 2, 'MedBay'),
       ('Woda', 'Butelki wody pitnej', 50, 'litry', 5, 'Storage Room'),
       ('Klucz francuski', 'Narzędzie do napraw', 5, 'szt.', 1, 'Workshop'),
       ('Żywność liofilizowana', 'Pakiety żywnościowe', 100, 'szt.', 10, 'Food Storage'),
       ('Kombinezon kosmiczny', 'Ochrona przed warunkami kosmicznymi', 3, 'szt.', 1, 'Equipment Room');