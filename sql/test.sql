CREATE TABLE credentials (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at timestamp,
    updated_at timestamp
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    username VARCHAR(100) NOT NULL,
    created_at timestamp,
    updated_at timestamp,
    FOREIGN KEY (username) REFERENCES credentials (username) ON DELETE CASCADE
);
