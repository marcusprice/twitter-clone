INSERT INTO Comment 
    (user_id, post_id, content, image, depth)
VALUES
    ($1, $2, $3, $4, 0)
RETURNING id;
