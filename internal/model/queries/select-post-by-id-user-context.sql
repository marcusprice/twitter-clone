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
    Post.updated_at,
    CASE
        WHEN PostLike.post_id IS NOT NULL THEN 1
        ELSE 0
    END AS liked
FROM
    Post
    INNER JOIN User ON User.id = Post.user_id
    LEFT JOIN PostLike ON PostLike.post_id = Post.id AND PostLike.user_id = $1
WHERE Post.id = $2;
