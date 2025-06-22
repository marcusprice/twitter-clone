SELECT
    Comment.id,
    Comment.post_id,
    Comment.user_id,
    Comment.depth,
    Comment.parent_comment_id,
    Comment.content,
    Comment.image,
    Comment.like_count,
    Comment.retweet_count,
    Comment.bookmark_count,
    Comment.impressions,
    Comment.created_at,
    Comment.updated_at,
    Author.user_name,
    Author.display_name,
    Author.avatar
FROM 
    Comment
    INNER JOIN User Author ON Author.id = Comment.user_id
WHERE Comment.id = $1;
