SELECT
	User.user_name,
	User.display_name,
	User.avatar
FROM
	User
WHERE 
    NOT EXISTS( 
        SELECT 1
        FROM UserFollows
        WHERE UserFollows.followee_id = User.id AND UserFollows.follower_id = $1
    )
    AND 
    User.id != $1
LIMIT 10;
