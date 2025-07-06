SELECT
	'post' AS type,
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
	'post-retweet' AS type,
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
UNION ALL
SELECT
	'comment-retweet' AS type,
    Comment.id,
    Comment.user_id,
    Comment.content,
    0 AS comment_count,
    Comment.like_count,
    Comment.retweet_count,
    Comment.bookmark_count,
    Comment.impressions,
    Comment.image,
    Comment.created_at,
    Comment.updated_at,
    Author.user_name,
    Author.display_name,
    Author.avatar,
    Retweeter.user_name AS retweeter_user_name,
    Retweeter.display_name AS retweeter_display_name,
    CASE
        WHEN CommentLike.comment_id IS NOT NULL THEN 1
        ELSE 0
    END AS liked,
    CASE
        WHEN ViewerRetweet.comment_id IS NOT NULL THEN 1
    	ELSE 0
    END AS retweeted,
    CASE
        WHEN CommentBookmark.comment_id IS NOT NULL THEN 1
    	ELSE 0
    END AS bookmarked,
    CommentRetweet.created_at AS sort_time
FROM
	CommentRetweet
	INNER JOIN Comment ON Comment.id = CommentRetweet.comment_id
	INNER JOIN User Retweeter ON Retweeter.id = CommentRetweet.user_id
	INNER JOIN User Author ON Author.id = Comment.user_id
	LEFT JOIN CommentLike ON CommentLike.comment_id = Comment.id AND CommentLike.user_id = $1
	LEFT JOIN CommentRetweet ViewerRetweet ON ViewerRetweet.comment_id = Comment.id AND ViewerRetweet.user_id = $1
	LEFT JOIN CommentBookmark ON CommentBookmark.comment_id = Comment.id AND CommentBookmark.user_id = $1
ORDER BY sort_time DESC
LIMIT $2 OFFSET $3;
