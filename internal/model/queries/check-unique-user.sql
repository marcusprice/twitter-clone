SELECT COUNT(*) FROM User
WHERE email = $1 OR user_name = $2;
