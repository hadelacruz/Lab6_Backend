-- Tabla de partidos
CREATE TABLE matches (
    id SERIAL PRIMARY KEY,
    homeTeam VARCHAR(255) NOT NULL,
    awayTeam VARCHAR(255) NOT NULL,
    matchDate DATE NOT NULL
);

ALTER TABLE matches
ADD COLUMN goals INT DEFAULT 0;
ADD COLUMN yellowCards INT DEFAULT 0;
ADD COLUMN redCards INT DEFAULT 0;
ADD COLUMN extratime INT DEFAULT 0;


