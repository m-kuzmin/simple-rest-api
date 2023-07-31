-- name: CreateUser :exec
INSERT INTO users (
    id, name, phone_number, country, city
) VALUES (
    $1, $2, $3, $4, $5
);

-- name: DeleteUserByID :exec
DELETE FROM users
WHERE id = $1;
