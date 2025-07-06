SELECT
    Post.id AS post_id,
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
    NULL AS retweeter_user_name,
    NULL AS retweeter_display_name,
    CASE
        WHEN PostLike.post_id IS NOT NULL THEN 1
        ELSE 0
    END AS liked,
    CASE
        WHEN ViewerRetweet.post_id IS NOT NULL THEN 1
    	ELSE 0
    END AS retweeted,
    CASE
        WHEN PostBookmark.post_id IS NOT NULL THEN 1
    	ELSE 0
    END AS bookmarked,
    Post.created_at AS sort_time
FROM 
    Post
    INNER JOIN User Author ON Author.id = Post.user_id
    LEFT JOIN PostLike ON PostLike.post_id = Post.id AND PostLike.user_id = $1
    LEFT JOIN PostRetweet ViewerRetweet ON ViewerRetweet.post_id = Post.id AND ViewerRetweet.user_id = $1
    LEFT JOIN PostBookmark ON PostBookmark.post_id = Post.id AND PostBookmark.user_id = $1
UNION ALL
SELECT
    Post.id AS post_id,
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
    Retweeter.user_name AS retweeter_user_name,
    Retweeter.display_name AS retweeter_display_name,
    CASE
        WHEN PostLike.post_id IS NOT NULL THEN 1
        ELSE 0
    END AS liked,
    CASE
        WHEN ViewerRetweet.post_id IS NOT NULL THEN 1
    	ELSE 0
    END AS retweeted,
    CASE
        WHEN PostBookmark.post_id IS NOT NULL THEN 1
    	ELSE 0
    END AS bookmarked,
    PostRetweet.created_at AS sort_time
FROM 
    PostRetweet
    INNER JOIN Post ON Post.id = PostRetweet.post_id
    INNER JOIN User Author ON Author.id = Post.user_id
    INNER JOIN User Retweeter ON Retweeter.id = PostRetweet.user_id
	LEFT JOIN PostLike ON PostLike.post_id = Post.id AND PostLike.user_id = $1
	LEFT JOIN PostRetweet ViewerRetweet ON ViewerRetweet.post_id = Post.id AND ViewerRetweet.user_id = $1
	LEFT JOIN PostBookmark ON PostBookmark.post_id = Post.id AND PostBookmark.user_id = $1
ORDER BY sort_time DESC
LIMIT $2 OFFSET $3;

