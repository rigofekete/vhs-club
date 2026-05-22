-- name: CreateRental :one
WITH new_rental AS (
  INSERT INTO rentals (user_id, tape_id)
  VALUES ($1, $2)
  RETURNING *
)
SELECT
  new_rental.*,
  tapes.title,
  users.username
FROM new_rental
JOIN tapes ON new_rental.tape_id = tapes.id
JOIN users ON new_rental.user_id = users.id;


INSERT INTO rentals (user_id, tape_id)
VALUES(
  $1,
  $2
)
RETURNING *;

-- name: ReturnTape :exec
UPDATE rentals
SET returned_at = NOW()
WHERE id = $1;

-- name: GetActiveRental :one
SELECT * FROM rentals
WHERE public_id = $1 AND user_id = $2 AND returned_at IS NULL;

-- name: GetActiveRentalsByUser :many
SELECT * FROM rentals
-- NULL is not a value so only IS keyword works
WHERE user_id = $1 AND returned_at IS NULL;

-- name: GetActiveRentalCountByUser :one
SELECT COUNT(*) FROM rentals
WHERE user_id = $1 AND returned_at IS NULL;

-- name: GetActiveRentalbyTape :many
SELECT * FROM rentals
WHERE tape_id = $1 AND returned_at IS NULL;

-- name: GetActiveRentalCountByTape :one
SELECT COUNT(*) FROM rentals
WHERE tape_id = $1 AND returned_at IS NULL;

-- name: GetUserActiveRentals :many
SELECT
  rentals.*,
  tapes.title,
  users.username
FROM rentals
JOIN tapes ON rentals.tape_id = tapes.id
JOIN users ON rentals.user_id = users.id
WHERE rentals.returned_at IS NULL AND rentals.user_id = $1
ORDER BY rentals.created_at ASC;

-- name: GetAllActiveRentals :many
SELECT
  rentals.*,
  tapes.title,
  users.username
FROM rentals
JOIN tapes ON rentals.tape_id = tapes.id
JOIN users ON rentals.user_id = users.id
WHERE returned_at IS NULL
ORDER BY rentals.created_at ASC;


-- name: DeleteAllRentals :exec
DELETE FROM rentals;
