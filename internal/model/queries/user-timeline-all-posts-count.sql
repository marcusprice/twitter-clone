SELECT (
    SELECT COUNT(*) FROM Post
) 
+ 
(
    SELECT COUNT(*) FROM PostRetweet
) 
+
(
    SELECT COUNT(*) FROM CommentRetweet
)
AS total_Count;
