INSERT INTO USER (email, user_name, password, first_name, last_name, display_name)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;
