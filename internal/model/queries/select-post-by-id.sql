SELECT
    User.user_name,
    User.display_name,
    User.avatar,
    Post.id,
    Post.user_id,
    Post.content,
    Post.comment_count,
    Post.like_count,
    Post.retweet_count,
    Post.bookmark_count,
    Post.impressions,
    Post.image,
    Post.created_at,
    Post.updated_at
FROM
    Post
    INNER JOIN User ON User.id = Post.user_id
WHERE Post.id = $1;
