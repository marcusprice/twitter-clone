UPDATE User
SET 
    last_login = CURRENT_TIMESTAMP,
    is_active = 1
WHERE id = $1
RETURNING
    last_login, is_active;
