-- Tabla de partidos
CREATE TABLE matches (
    id SERIAL PRIMARY KEY,
    homeTeam VARCHAR(255) NOT NULL,
    awayTeam VARCHAR(255) NOT NULL,
    matchDate DATE NOT NULL
);
