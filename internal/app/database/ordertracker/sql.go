package ordertracker

const acquireSQL = `
UPDATE Queue SET lock = TRUE
WHERE order_id = (
	SELECT order_id
	FROM Queue
	WHERE status_id IN (0, 1, 2) AND lock = FALSE
	ORDER BY updated_at
	LIMIT 1
	FOR UPDATE SKIP LOCKED
)
RETURNING order_id, status_id;
`

const updateAndReleaseSQL = `
UPDATE Queue
SET lock = FALSE, status_id = $1, updated_at = CURRENT_TIMESTAMP
WHERE order_id=$2;
`

const addSQL = `
INSERT INTO Queue(order_id) VALUES ($1);
`

const deleteSQL = `
DELETE FROM Queue WHERE order_id = $1;
`
