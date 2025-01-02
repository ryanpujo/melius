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

-- Country Table
CREATE TABLE countries (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

-- State Table
CREATE TABLE states (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    country_id INT NOT NULL REFERENCES countries(id) ON DELETE CASCADE
);

-- City Table
CREATE TABLE cities (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    state_id INT NOT NULL REFERENCES states(id) ON DELETE CASCADE
);

-- Address Table
CREATE TABLE addresses (
    id SERIAL PRIMARY KEY,
    address_line TEXT NOT NULL,
    postal_code VARCHAR(50) NOT NULL,
    is_main BOOLEAN NOT NULL,
    city_id INT NOT NULL REFERENCES cities(id) ON DELETE CASCADE
);

CREATE TABLE user_address (
    address_id INT NOT NULL REFERENCES addresses(id) ON DELETE CASCADE
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE
);
