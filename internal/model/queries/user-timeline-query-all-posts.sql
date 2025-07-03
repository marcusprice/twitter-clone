SELECT
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
    Author.user_name,
    Author.display_name,
    Author.avatar,
    Retweeter.user_name,
    Retweeter.display_name,
    CASE
        WHEN PostLike.post_id IS NOT NULL THEN 1
        ELSE 0
    END AS liked
FROM 
    Post
    INNER JOIN User Author ON Author.id = Post.user_id
    LEFT JOIN PostRetweet ON PostRetweet.post_id = Post.id
    LEFT JOIN User Retweeter ON Retweeter.id = PostRetweet.user_id
    LEFT JOIN PostLike ON PostLike.post_id = Post.id AND PostLike.user_id = $1
ORDER BY 
    COALESCE(PostRetweet.created_at, Post.created_at) DESC
LIMIT $2 OFFSET $3;
