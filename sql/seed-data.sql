-- Users
INSERT INTO User (id, email, user_name, password, first_name, last_name, display_name)
VALUES
  (1, 'alice@example.com', 'alice', 'hashed_pw_1', 'Alice', 'Anderson', 'alice'),
  (2, 'bob@example.com', 'bob', 'hashed_pw_2', 'Bob', 'Baxter', 'bobby'),
  (3, 'carol@example.com', 'carol', 'hashed_pw_3', 'Carol', 'Clark', 'carolC'),
  (4, 'dave@example.com', 'dave', 'hashed_pw_4', 'Dave', 'Dixon', 'daveD'),
  (5, 'eve@example.com', 'eve', 'hashed_pw_5', 'Eve', 'Evans', 'eveE');

-- Follows
INSERT INTO UserFollows (follower_id, followee_id)
VALUES
  (1, 2), (1, 3), (2, 3), (3, 1), (4, 1);

-- Posts
INSERT INTO Post (id, user_id, content, image_url)
VALUES
  (1, 2, 'Hello world from Bob!', NULL),
  (2, 3, 'Carol here with a picture 📸', 'https://example.com/pic1.jpg'),
  (3, 1, 'Alice posting something interesting.', NULL);

-- Post Likes
INSERT INTO PostLike (post_id, user_id) VALUES
  (1, 1), (1, 3), (2, 1), (2, 2);

-- Post Retweets
INSERT INTO PostRetweet (post_id, user_id) VALUES
  (2, 4), (3, 2);

-- Post Bookmarks
INSERT INTO PostBookmark (post_id, user_id) VALUES
  (3, 5), (1, 4);

-- Comments
INSERT INTO Comment (id, content, post_id, user_id)
VALUES
  (1, 'Nice post Bob!', 1, 1),
  (2, 'Thanks Alice!', 1, 2),
  (3, 'Agree with Carol!', 2, 5);

INSERT INTO Comment (id, content, post_id, user_id, parent_comment_id)
VALUES
  (4, 'Reply to a comment', 1, 3, 1);

-- Comment Likes
INSERT INTO CommentLike (comment_id, user_id) VALUES
  (1, 2), (2, 1);

-- Comment Retweets
INSERT INTO CommentRetweet (comment_id, user_id) VALUES
  (1, 3);

-- Comment Bookmarks
INSERT INTO CommentBookmark (comment_id, user_id) VALUES
  (2, 5);

-- Notifications
INSERT INTO Notification (initiator_id, receiver_id, post_id, type, is_read)
VALUES
  (1, 2, 1, 'post_like', 0),
  (3, 2, 1, 'post_like', 1),
  (5, 3, 2, 'post_comment', 0),
  (2, 1, 3, 'post_retweet', 0);
