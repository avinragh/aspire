CREATE USER avinragh WITH PASSWORD 'toor';
CREATE DATABASE aspire OWNER avinragh;
GRANT ALL PRIVILEGES ON DATABASE aspire TO avinragh;
\connect aspire avinragh;
BEGIN;
CREATE TABLE if NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email varchar(225) not null unique,
    username varchar(225),
    password varchar(225) not null,
    role varchar(20),
    created_on timestamp,
    modified_on timestamp
);

CREATE TABLE IF NOT EXISTS loans (
    id SERIAL PRIMARY KEY,
    amount DECIMAL NOT NULL,
    currency VARCHAR(4),
    term INT NOT NULL,
    state VARCHAR(10),
    start_date TIMESTAMP,
    created_on TIMESTAMP,
    modified_on TIMESTAMP,
    user_id INT,
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
            REFERENCES users(id)

);
CREATE TABLE IF NOT EXISTS installments (
    id SERIAL PRIMARY KEY,
    installment_amount DECIMAL NOT NULL,
    due_date TIMESTAMP,
    state VARCHAR(10),
    loan_id INT,
    created_on TIMESTAMP,
    modified_on TIMESTAMP,
    CONSTRAINT fk_loan
        FOREIGN KEY(loan_id)
            REFERENCES loans(id)

);
COMMIT;
