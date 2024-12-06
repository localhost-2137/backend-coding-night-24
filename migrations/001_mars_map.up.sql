CREATE TABLE mars_map_point (
    id INTEGER PRIMARY KEY,
    longitude DECIMAL(9, 6) NOT NULL,
    latitude DECIMAL(9, 6) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    label TEXT NOT NULL
);
