INSERT INTO Post (user_id, content, image)
VALUES ($1, $2, $3)
RETURNING id;
