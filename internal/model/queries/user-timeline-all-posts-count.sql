SELECT (
    SELECT COUNT(*) FROM Post
) 
+ 
(
    SELECT COUNT(*) FROM PostRetweet
) AS total_Count;
