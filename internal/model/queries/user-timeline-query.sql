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
    Author.user_name,
    Author.display_name,
    Author.avatar,
    Retweeter.user_name,
    Retweeter.display_name
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
    COALESCE(PostRetweet.created_at, Post.created_at) DESC
LIMIT $2 OFFSET $3;
