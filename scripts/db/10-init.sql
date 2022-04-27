CREATE USER avinragh WITH PASSWORD 'toor';
CREATE DATABASE aspire OWNER avinragh;
  GRANT ALL PRIVILEGES ON DATABASE aspire TO avinragh;
  \connect aspire avinragh
  BEGIN;
    CREATE TABLE IF NOT EXISTS loan (
	  id SERIAL PRIMARY KEY,
	  amount DECIMAL NOT NULL,
	  currency VARCHAR(4),
      term INT NOT NULL,
      state VARCHAR(10),
      starDate TIMESTAMP,
      createdOn TIMESTAMP,
      modifiedOn TIMESTAMP,
      userId INT,
	);
    CREATE TABLE IF NOT EXISTS installment (
        id SERIAL PRIMARY KEY,
        installmentAmount DECIMAL NOT NULL,
        dueDate TIMESTAMP,
        state VARCHAR(10),
        loanId INT,
        createdOn TIMESTAMP,
        modifiedOn TIMESTAMP,
        CONSTRAINT fk_loan
            FOREIGN KEY(installmentId)
                REFERENCES loan(installmentId)

    );
  COMMIT;
EOSQL
      installmentAmount:
        type: integer
        format: float64
      dueDate:
        type: string
        format: date-time
      state:
        type: string
      installmentId:
        type: integer
        format: int64
      createdOn:
        type: string
        format: date-time
      modifiedOn:
        type: string
        format: date-time
