SELECT
    Author.user_name AS username,
    Author.display_name AS display_name,
    Author.avatar AS avatar,
    Author.bio AS bio,
    COUNT(Followers.id) AS follower_count,
    COUNT(Following.id) AS following_count,
    CASE WHEN ViewerFollowing.id IS NOT NULL THEN 1 ELSE 0 END AS viewer_following
FROM
    Post
    INNER JOIN User Author
        ON Author.id = Post.user_id
    LEFT JOIN UserFollows Followers
        ON Followers.followee_id = Author.id
    LEFT JOIN UserFollows Following
        ON Following.follower_id = Author.id
    LEFT JOIN UserFollows ViewerFollowing
        ON ViewerFollowing.followee_id = Author.id AND ViewerFollowing.follower_id = $1
WHERE Post.id = $2
GROUP BY
    Author.id;
