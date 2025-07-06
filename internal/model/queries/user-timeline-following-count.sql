SELECT(
	SELECT COUNT(*)
	FROM
		Post
		INNER JOIN User Author ON Author.id = Post.user_id
		LEFT JOIN UserFollows ON UserFollows.followee_id = Post.user_id AND UserFollows.follower_id = $1
	WHERE UserFollows.follower_id IS NOT NULL
)
+
(
	SELECT COUNT(*)
	FROM
		PostRetweet
		INNER JOIN Post ON Post.id = PostRetweet.post_id
		INNER JOIN User Retweeter ON Retweeter.id = PostRetweet.user_id
		INNER JOIN User Author ON Author.id = Post.user_id
		LEFT JOIN UserFollows ON UserFollows.followee_id = PostRetweet.user_id
	WHERE UserFollows.follower_id = $1
) AS total_count;
