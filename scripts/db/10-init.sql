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
    createdOn timestamp,
    modifiedOn timestamp
);

CREATE TABLE IF NOT EXISTS loans (
    id SERIAL PRIMARY KEY,
    amount DECIMAL NOT NULL,
    currency VARCHAR(4),
    term INT NOT NULL,
    state VARCHAR(10),
    starDate TIMESTAMP,
    createdOn TIMESTAMP,
    modifiedOn TIMESTAMP,
    userId INT,
    CONSTRAINT fk_user
        FOREIGN KEY(userId)
            REFERENCES users(id)

);
CREATE TABLE IF NOT EXISTS installments (
    id SERIAL PRIMARY KEY,
    installmentAmount DECIMAL NOT NULL,
    dueDate TIMESTAMP,
    state VARCHAR(10),
    loanId INT,
    createdOn TIMESTAMP,
    modifiedOn TIMESTAMP,
    CONSTRAINT fk_loan
        FOREIGN KEY(loanId)
            REFERENCES loans(id)

);
COMMIT;
