-- name: CreateUser :exec
INSERT INTO users (
    id, name, phone_number, country, city
) VALUES (
    $1, $2, $3, $4, $5
);

-- name: DeleteUserByID :exec
DELETE FROM users
WHERE id = $1;

-- name: SearchUsers :many
SELECT *
FROM users
WHERE
    (name         LIKE '%' || $1 || '%' OR $1 IS NULL) AND
    (phone_number LIKE '%' || $2 || '%' OR $2 IS NULL) AND
    (country      LIKE '%' || $3 || '%' OR $3 IS NULL) AND
    (city         LIKE '%' || $4 || '%' OR $4 IS NULL);
