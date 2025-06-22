INSERT INTO Comment 
    (user_id, post_id, parent_comment_id, content, image, depth)
VALUES
    ($1, $2, $3, $4, $5, 1)
RETURNING id;
