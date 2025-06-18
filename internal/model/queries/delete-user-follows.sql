DELETE FROM UserFollows WHERE follower_id = $1 AND followee_id = $2;
