UPDATE User
SET last_login = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING last_login;
