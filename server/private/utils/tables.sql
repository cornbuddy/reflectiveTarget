CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(64) NOT NULL UNIQUE,
    hashed_password VARCHAR(128) NOT NULL
);

CREATE TABLE IF NOT EXISTS targets (
    id SERIAL PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    owner_id INT REFERENCES users(id) NOT NULL
);

CREATE TABLE IF NOT EXISTS questions (
    id SERIAL PRIMARY KEY NOT NULL,
    text VARCHAR(128) NOT NULL,
    target_id INT REFERENCES targets(id)
);

CREATE TABLE IF NOT EXISTS shots (
    id SERIAL PRIMARY KEY,
    x INT NOT NULL,
    y INT NOT NULL,
    target_id INT REFERENCES targets(id) NOT NULL,
    shooter_id INT REFERENCES users(id) NOT NULL
);
