CREATE TABLE movies (
    id BIGSERIAL PRIMARY KEY,
    prefix text NOT NULL,
    title text NOT NULL,
    description text
);
