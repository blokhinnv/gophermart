package database

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
const getOrdersByUserID = `
WITH a AS (
    SELECT order_id, SUM(sum) as sum
    FROM Transaction
    WHERE transaction_type_id=1
    GROUP BY order_id
)
SELECT o.id AS "number", s.status, a.sum AS "accrual", o.uploaded_at
FROM UserOrder o
LEFT JOIN a ON o.id = a.order_id
JOIN OrderStatus s ON s.id = o.status_id
WHERE o.user_id = $1
ORDER BY o.uploaded_at;
`
const getBalanceSQL = `
SELECT (
	COALESCE(SUM(CASE WHEN transaction_type_id=1 THEN t.sum END), 0) -
	COALESCE(SUM(CASE WHEN transaction_type_id=2 THEN t.sum END), 0)
) AS balance, (
	COALESCE(SUM(CASE WHEN transaction_type_id=2 THEN t.sum END), 0)
) AS withdrawn
FROM Transaction t
JOIN UserOrder o ON o.id = t.order_id
WHERE o.user_id = $1;
`

const getWithdrawalsSQL = `
SELECT order_id AS order, sum, processed_at
FROM Transaction t
JOIN UserOrder o ON o.id = t.order_id
WHERE o.user_id = $1 AND t.transaction_type_id=2
ORDER BY t.processed_at;
`
