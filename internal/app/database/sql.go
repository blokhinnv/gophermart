package database

const createSQL = `
DROP TABLE IF EXISTS UserAccount CASCADE;
CREATE TABLE UserAccount(
	id SERIAL PRIMARY KEY,
	username VARCHAR UNIQUE NOT NULL,
	hashed_password VARCHAR NOT NULL,
	salt VARCHAR NOT NULL
);

DROP TABLE IF EXISTS OrderStatus CASCADE;
CREATE TABLE OrderStatus(
	id INTEGER PRIMARY KEY,
	status VARCHAR NOT NULL
);
INSERT INTO OrderStatus VALUES
(0, 'NEW'), (1, 'REGISTERED'), (2, 'PROCESSING'), (3, 'INVALID'), (4, 'PROCESSED');

DROP TABLE IF EXISTS UserOrder CASCADE;
CREATE TABLE UserOrder(
	id VARCHAR PRIMARY KEY,
	user_id INTEGER NOT NULL,
	status_id INTEGER NOT NULL,
	updateded_at TIMESTAMP DEFAULT NOW(),
	CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES UserAccount(id),
	CONSTRAINT fk_status_id FOREIGN KEY (status_id) REFERENCES OrderStatus(id)
);

DROP TABLE IF EXISTS TransactionType CASCADE;
CREATE TABLE TransactionType(
	id SERIAL PRIMARY KEY,
	type VARCHAR NOT NULL
);
INSERT INTO TransactionType(type) VALUES
('Списание баллов лояльности'), ('Зачисление баллов лояльности');

DROP TABLE IF EXISTS Transaction CASCADE;
CREATE TABLE Transaction(
	id SERIAL PRIMARY KEY,
	order_id VARCHAR NOT NULL,
	sum DOUBLE PRECISION NOT NULL,
	transaction_type_id INTEGER NOT NULL,
	processed_at TIMESTAMP DEFAULT NOW(),
	CONSTRAINT fk_transaction_type_id
		FOREIGN KEY (transaction_type_id) REFERENCES TransactionType(id),
	CONSTRAINT fk_order_id FOREIGN KEY (order_id) REFERENCES UserOrder(id)
);
`

const addUserSQL = `
INSERT INTO UserAccount(username, hashed_password, salt) VALUES ($1, $2, $3);
`
const selectUserByLoginSQL = `
SELECT hashed_password, salt FROM UserAccount WHERE username=$1;
`
