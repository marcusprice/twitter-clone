SELECT
    COUNT(*)
FROM 
    Post
    INNER JOIN User Author ON Author.id = Post.user_id
    LEFT JOIN PostRetweet ON PostRetweet.post_id = Post.id
    LEFT JOIN User PostRetweeter ON PostRetweeter.id = PostRetweet.user_id;
