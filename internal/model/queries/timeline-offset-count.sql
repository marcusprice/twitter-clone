SELECT
    COUNT(*)
FROM
    Post
    INNER JOIN User ON User.id = Post.user_id
    INNER JOIN UserFollows ON UserFollows.followee_id = User.id
WHERE
    UserFollows.follower_id = $1
ORDER BY
    Post.created_at DESC;
