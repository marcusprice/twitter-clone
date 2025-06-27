SELECT
    PostBookmark.created_at AS bookmark_created_at,
    Post.id AS id,
    Post.content AS content,
    Post.image AS image,
    Post.like_count AS like_count,
    Post.retweet_count AS retweet_count,
    Post.bookmark_count AS bookmark_count,
    Post.impressions AS impressions,
    Post.created_at AS created_at,
    Post.updated_at AS updated_at,
    PostAuthor.user_name as author_user_name,
    PostAuthor.display_name as author_display_name,
    PostAuthor.avatar as author_avatar,
    'post' AS type
FROM
    PostBookmark
    INNER JOIN Post ON Post.id = PostBookmark.post_id
    INNER JOIN User PostAuthor ON PostAuthor.id = Post.user_id

WHERE PostBookmark.user_id = $1

UNION

SELECT
    CommentBookmark.created_at AS bookmark_created_at,
    Comment.id AS id,
    Comment.content AS content,
    Comment.image AS image,
    Comment.like_count AS like_count,
    Comment.retweet_count AS retweet_count,
    Comment.bookmark_count AS bookmark_count,
    Comment.impressions AS impressions,
    Comment.created_at AS created_at,
    Comment.updated_at AS updated_at,
    CommentAuthor.user_name as author_user_name,
    CommentAuthor.display_name as author_display_name,
    CommentAuthor.avatar as author_avatar,
    'comment' AS type
FROM
    CommentBookmark
    INNER JOIN Comment ON Comment.id = CommentBookmark.comment_id
    INNER JOIN User CommentAuthor ON CommentAuthor.id = Comment.user_id

WHERE CommentBookmark.user_id = $1

ORDER BY
    bookmark_created_at DESC
LIMIT $2 OFFSET $3;
