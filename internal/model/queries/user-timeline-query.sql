SELECT
    Post.id,
    Post.user_id,
    Post.content,
    Post.like_count,
    Post.retweet_count,
    Post.bookmark_count,
    Post.impressions,
    Post.image,
    Post.created_at,
    Post.updated_at,
    User.user_name,
    User.display_name,
    User.avatar
FROM
    Post
    INNER JOIN User ON User.id = Post.user_id
    INNER JOIN UserFollows ON UserFollows.followee_id = User.id
WHERE
    UserFollows.follower_id = $1
ORDER BY 
    Post.created_at DESC
LIMIT $2 OFFSET $3;
