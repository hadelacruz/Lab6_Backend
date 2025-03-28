CREATE TABLE IF NOT EXISTS matches (
    id SERIAL PRIMARY KEY,
    team1 VARCHAR(100) NOT NULL,
    team2 VARCHAR(100) NOT NULL,
    score VARCHAR(50)
);

-- Insertar algunos partidos de prueba
INSERT INTO matches (team1, team2, score) VALUES
    ('Team A', 'Team B', '2-1'),
    ('Team C', 'Team D', '3-0');
