DROP TABLE IF EXISTS album;
CREATE TABLE album (
    id SERIAL PRIMARY KEY,
    title VARCHAR(128) NOT NULL,
    artist VARCHAR(255) NOT NULL,
    price REAL NOT NULL
);

INSERT INTO album 
    (title, artist, price)
VALUES 
	('Metallica (Black Album)', 'Metallica', 60.0),
    ('Piece Sells... but Whos Buying?', 'Megadeth', 59.99),
    ('Cowboys From Hell', 'Pantera',  75.0);