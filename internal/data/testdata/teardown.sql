DROP TABLE IF EXISTS movies;

ALTER TABLE movies DROP CONSTRAINT IF EXISTS movies_runtime_check;

ALTER TABLE movies DROP CONSTRAINT IF EXISTS movies_year_check;

ALTER TABLE movies DROP CONSTRAINT IF EXISTS genres_length_check;

DROP TABLE IF EXISTS users;

DROP TABLE IF EXISTS tokens;

DROP INDEX IF EXISTS movies_title_idx;
DROP INDEX IF EXISTS movies_genres_idx;

DROP TABLE IF EXISTS users_permissions;
DROP TABLE IF EXISTS permissions;