SELECT
    COUNT(*)
FROM
    Post
    INNER JOIN User Author ON Author.id = Post.user_id
    LEFT JOIN PostRetweet ON PostRetweet.post_id = Post.id
    LEFT JOIN User Retweeter ON Retweeter.id = PostRetweet.user_id
    LEFT JOIN UserFollows FollowedUsers ON FollowedUsers.followee_id = Author.id AND FollowedUsers.follower_id = $1
    LEFT JOIN UserFollows RetweetedUsers ON RetweetedUsers.followee_id = Retweeter.id AND RetweetedUsers.follower_id = $1
WHERE
    FollowedUsers.follower_id IS NOT NULL OR RetweetedUsers.follower_id IS NOT NULL
ORDER BY 
    Post.created_at DESC
