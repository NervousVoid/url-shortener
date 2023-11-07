DROP TABLE IF EXISTS urls;
CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    full_url VARCHAR(255) NOT NULL UNIQUE,
    short_url VARCHAR(20) NOT NULL UNIQUE
);