SELECT (
    SELECT COUNT(*)
    FROM PostBookmark
    INNER JOIN Post ON Post.id = PostBookmark.post_id
    INNER JOIN User PostAuthor ON PostAuthor.id = Post.user_id
    WHERE PostBookmark.user_id = $1
) +
(
    SELECT COUNT(*)
    FROM CommentBookmark
    INNER JOIN Comment ON Comment.id = CommentBookmark.comment_id
    INNER JOIN User CommentAuthor ON CommentAuthor.id = Comment.user_id
    WHERE CommentBookmark.user_id = $1
) AS total_count;
