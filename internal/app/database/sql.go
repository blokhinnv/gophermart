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
	status_id INTEGER NOT NULL DEFAULT 0,
	uploaded_at TIMESTAMP DEFAULT NOW(),
	CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES UserAccount(id),
	CONSTRAINT fk_status_id FOREIGN KEY (status_id) REFERENCES OrderStatus(id)
);

DROP TABLE IF EXISTS TransactionType CASCADE;
CREATE TABLE TransactionType(
	id SERIAL PRIMARY KEY,
	type VARCHAR NOT NULL
);
INSERT INTO TransactionType(type) VALUES
('ACCRUAL'), ('WITHDRAWAL');

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
INSERT INTO UserAccount(username, hashed_password, salt) VALUES ($1, $2, $3) RETURNING id;
`
const selectUserByLoginSQL = `
SELECT id, hashed_password, salt FROM UserAccount WHERE username=$1;
`
const addOrderSQL = `
INSERT INTO UserOrder(id, user_id) VALUES ($1, $2);
`
const selectOrderByIDSQL = `
SELECT id, user_id, status_id, uploaded_at FROM UserOrder WHERE id=$1;
`
const updateOrderStatusSQL = `
UPDATE UserOrder SET status_id=(SELECT id FROM OrderStatus WHERE status=$1) WHERE id=$2;
`
const addTransactionSQL = `
INSERT INTO Transaction(order_id, sum, transaction_type_id)
	SELECT $1, $2, id
	FROM TransactionType
	WHERE type=$3;
`
