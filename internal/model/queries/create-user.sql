INSERT INTO USER (email, user_name, password, first_name, last_name, display_name, is_active)
VALUES ($1, $2, $3, $4, $5, $6, 0)
RETURNING id, last_login, created_at, updated_at;
