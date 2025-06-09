SELECT 
    id,
    email,
    user_name,
    first_name,
    last_name,
    display_name,
    last_login,
    created_at,
    updated_at
FROM USER 
WHERE id = $1;
