
CREATE TABLE UserAccount(
	id SERIAL PRIMARY KEY,
	username VARCHAR UNIQUE NOT NULL,
	hashed_password VARCHAR NOT NULL,
	salt VARCHAR NOT NULL
);


CREATE TABLE OrderStatus(
	id INTEGER PRIMARY KEY,
	status VARCHAR NOT NULL
);
INSERT INTO OrderStatus VALUES
(0, 'NEW'), (1, 'REGISTERED'), (2, 'PROCESSING'), (3, 'INVALID'), (4, 'PROCESSED');


CREATE TABLE UserOrder(
	id VARCHAR PRIMARY KEY,
	user_id INTEGER NOT NULL,
	status_id INTEGER NOT NULL DEFAULT 0,
	uploaded_at TIMESTAMP DEFAULT NOW(),
	CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES UserAccount(id),
	CONSTRAINT fk_status_id FOREIGN KEY (status_id) REFERENCES OrderStatus(id)
);


CREATE TABLE TransactionType(
	id SERIAL PRIMARY KEY,
	type VARCHAR NOT NULL
);
INSERT INTO TransactionType(type) VALUES
('ACCRUAL'), ('WITHDRAWAL');


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


CREATE TABLE Queue(
	id SERIAL PRIMARY KEY,
	order_id VARCHAR NOT NULL,
	status_id INTEGER NOT NULL DEFAULT 0,
	lock BOOLEAN NOT NULL DEFAULT FALSE,
	updated_at TIMESTAMP DEFAULT NOW(),
	CONSTRAINT fk_status_id FOREIGN KEY (status_id) REFERENCES OrderStatus(id)
);
