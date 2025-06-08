DROP TRIGGER IF EXISTS update_user_timestamp;
DROP TRIGGER IF EXISTS update_post_timestamp;
DROP TRIGGER IF EXISTS update_comment_timestamp;
DROP TRIGGER IF EXISTS update_notification_timestamp;

DROP TRIGGER IF EXISTS increment_post_like_count;
DROP TRIGGER IF EXISTS decrement_post_like_count;
DROP TRIGGER IF EXISTS increment_post_retweet_count;
DROP TRIGGER IF EXISTS decrement_post_retweet_count;
DROP TRIGGER IF EXISTS increment_post_bookmark_count;
DROP TRIGGER IF EXISTS decrement_post_bookmark_count;

DROP TRIGGER IF EXISTS increment_comment_like_count;
DROP TRIGGER IF EXISTS decrement_comment_like_count;
DROP TRIGGER IF EXISTS increment_comment_retweet_count;
DROP TRIGGER IF EXISTS decrement_comment_retweet_count;
DROP TRIGGER IF EXISTS increment_comment_bookmark_count;
DROP TRIGGER IF EXISTS decrement_comment_bookmark_count;

DROP TABLE IF EXISTS User;
DROP TABLE IF EXISTS UserFollows;
DROP TABLE IF EXISTS Post;
DROP TABLE IF EXISTS Comment;
DROP TABLE IF EXISTS Notification;

DROP TABLE IF EXISTS PostLike;
DROP TABLE IF EXISTS PostRetweet;
DROP TABLE IF EXISTS PostBookmark;
DROP TABLE IF EXISTS CommentLike;
DROP TABLE IF EXISTS CommentRetweet;
DROP TABLE IF EXISTS CommentBookmarks;

CREATE TABLE User (
    id INTEGER PRIMARY KEY,
    email TEXT NOT NULL UNIQUE CHECK (length(trim(email)) > 0),
    user_name TEXT NOT NULL UNIQUE CHECK (length(trim(user_name)) > 0),
    password TEXT NOT NULL CHECK (length(trim(password)) > 0),
    first_name TEXT,
    last_name TEXT,
    display_name TEXT NOT NULL CHECK (length(trim(display_name)) > 0),
    created_at TEXT NOT NULL DEFAULT current_timestamp,
    updated_at TEXT NOT NULL DEFAULT current_timestamp
);

CREATE TABLE UserFollows (
    id INTEGER PRIMARY KEY,
    follower_id INTEGER NOT NULL,
    followee_id INTEGER NOT NULL,
    created_at TEXT NOT NULL DEFAULT current_timestamp,

    FOREIGN KEY (follower_id) REFERENCES User (id) ON DELETE CASCADE,
    FOREIGN KEY (followee_id) REFERENCES User (id) ON DELETE CASCADE,

    UNIQUE (follower_id, followee_id),
    CHECK (follower_id != followee_id)
);

CREATE TABLE Post (
    id INTEGER PRIMARY KEY,
    user_id INTEGER,
    content TEXT NOT NULL CHECK (length(trim(content)) > 0),
    like_count INTEGER DEFAULT 0,
    retweet_count INTEGER DEFAULT 0,
    bookmark_count INTEGER DEFAULT 0,
    impressions INTEGER DEFAULT 0,
    image_url TEXT,
    created_at TEXT NOT NULL DEFAULT current_timestamp,
    updated_at TEXT NOT NULL DEFAULT current_timestamp,

    FOREIGN KEY (user_id) REFERENCES User (id) ON DELETE CASCADE

    CHECK (like_count >= 0),
    CHECK (retweet_count >= 0),
    CHECK (bookmark_count >= 0)
);

CREATE TABLE PostLike (
    id INTEGER PRIMARY KEY,
    post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    created_at TEXT NOT NULL DEFAULT current_timestamp,

    FOREIGN KEY (post_id) REFERENCES Post (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES User (id) ON DELETE CASCADE
);

CREATE TABLE PostRetweet (
    id INTEGER PRIMARY KEY,
    content TEXT,
    image_url TEXT,
    post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    created_at TEXT NOT NULL DEFAULT current_timestamp,

    FOREIGN KEY (post_id) REFERENCES Post (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES User (id) ON DELETE CASCADE
);

CREATE TABLE PostBookmark (
    id INTEGER PRIMARY KEY,
    post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    created_at TEXT NOT NULL DEFAULT current_timestamp,

    FOREIGN KEY (post_id) REFERENCES Post (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES User (id) ON DELETE CASCADE
);

CREATE TABLE Comment (
    id INTEGER PRIMARY KEY,
    content TEXT,
    like_count INTEGER DEFAULT 0,
    retweet_count INTEGER DEFAULT 0,
    bookmark_count INTEGER DEFAULT 0,
    impressions INTEGER DEFAULT 0,
    post_id INTEGER,
    user_id INTEGER,
    parent_comment_id INTEGER,
    created_at TEXT NOT NULL DEFAULT current_timestamp,
    updated_at TEXT NOT NULL DEFAULT current_timestamp,

    FOREIGN KEY (post_id) REFERENCES Post (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES User (id) ON DELETE CASCADE,
    FOREIGN KEY (parent_comment_id) REFERENCES Comment (id) ON DELETE CASCADE

    CHECK (like_count >= 0),
    CHECK (retweet_count >= 0),
    CHECK (bookmark_count >= 0)
);

CREATE TABLE CommentLike (
    id INTEGER PRIMARY KEY,
    comment_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    created_at TEXT NOT NULL DEFAULT current_timestamp,

    FOREIGN KEY (comment_id) REFERENCES Comment (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES User (id) ON DELETE CASCADE
);

CREATE TABLE CommentRetweet (
    id INTEGER PRIMARY KEY,
    comment_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    created_at TEXT NOT NULL DEFAULT current_timestamp,

    FOREIGN KEY (comment_id) REFERENCES Comment (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES User (id) ON DELETE CASCADE
);

CREATE TABLE CommentBookmark (
    id INTEGER PRIMARY KEY,
    comment_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    created_at TEXT NOT NULL DEFAULT current_timestamp,

    FOREIGN KEY (comment_id) REFERENCES Comment (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES User (id) ON DELETE CASCADE
);

CREATE TABLE Notification (
    id INTEGER PRIMARY KEY,
    initiator_id INTEGER,
    receiver_id INTEGER,
    post_id INTEGER,
    comment_id INTEGER,
    type TEXT NOT NULL CHECK (type IN(
        'post_like', 'post_comment','post_retweet',
        'comment_like', 'comment_reply', 'comment_retweet', 
        'mention'
    )),
    is_read INTEGER NOT NULL CHECK (is_read IN(0, 1)),
    created_at TEXT NOT NULL DEFAULT current_timestamp,
    updated_at TEXT NOT NULL DEFAULT current_timestamp,

    FOREIGN KEY (initiator_id) REFERENCES User (id) ON DELETE CASCADE,
    FOREIGN KEY (receiver_id) REFERENCES User (id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES Post (id) ON DELETE CASCADE,
    FOREIGN KEY (comment_id) REFERENCES Comment (id) ON DELETE CASCADE,

    CHECK(
        (post_id IS NOT NULL AND comment_id IS NULL) OR
        (post_id IS NULL AND comment_id IS NOT NULL)
    )
);

CREATE TRIGGER update_user_timestamp
AFTER UPDATE ON User
BEGIN
    UPDATE User SET updated_at = current_timestamp WHERE id = NEW.id;
END;

CREATE TRIGGER update_post_timestamp
AFTER UPDATE ON Post
BEGIN
    UPDATE Post SET updated_at = current_timestamp WHERE id = NEW.id;
END;

CREATE TRIGGER update_comment_timestamp
AFTER UPDATE ON Comment
BEGIN
    UPDATE Comment SET updated_at = current_timestamp WHERE id = NEW.id;
END;

CREATE TRIGGER update_notification_timestamp
AFTER UPDATE ON Notification
BEGIN
    UPDATE Notification SET updated_at = current_timestamp WHERE id = NEW.id;
END;



CREATE TRIGGER increment_post_like_count
AFTER INSERT ON PostLike
BEGIN
    UPDATE Post SET like_count = like_count + 1 WHERE id = NEW.post_id;
END;

CREATE TRIGGER decrement_post_like_count
AFTER DELETE ON PostLike
BEGIN
    UPDATE Post SET like_count = like_count - 1 WHERE id = NEW.post_id;
END;

CREATE TRIGGER increment_post_retweet_count
AFTER INSERT ON PostRetweet
BEGIN
    UPDATE Post SET retweet_count = retweet_count + 1 WHERE id = NEW.post_id;
END;

CREATE TRIGGER decrement_post_retweet_count
AFTER DELETE ON PostRetweet
BEGIN
    UPDATE Post SET retweet_count = retweet_count - 1 WHERE id = NEW.post_id;
END;

CREATE TRIGGER increment_post_bookmark_count
AFTER INSERT ON PostBookmark
BEGIN
    UPDATE Post SET bookmark_count = bookmark_count + 1 WHERE id = NEW.post_id;
END;

CREATE TRIGGER decrement_post_bookmark_count
AFTER DELETE ON PostBookmark
BEGIN
    UPDATE Post SET bookmark_count = bookmark_count - 1 WHERE id = NEW.post_id;
END;

CREATE TRIGGER increment_comment_like_count
AFTER INSERT ON CommentLike
BEGIN
    UPDATE Comment SET like_count = like_count + 1 WHERE id = NEW.comment_id;
END;

CREATE TRIGGER decrement_comment_like_count
AFTER DELETE ON CommentLike
BEGIN
    UPDATE Comment SET like_count = like_count - 1 WHERE id = NEW.comment_id;
END;

CREATE TRIGGER increment_comment_retweet_count
AFTER INSERT ON CommentRetweet
BEGIN
    UPDATE Comment SET retweet_count = retweet_count + 1 WHERE id = NEW.comment_id;
END;

CREATE TRIGGER decrement_comment_retweet_count
AFTER DELETE ON CommentRetweet
BEGIN
    UPDATE Comment SET retweet_count = retweet_count - 1 WHERE id = NEW.comment_id;
END;

CREATE TRIGGER increment_comment_bookmark_count
AFTER INSERT ON CommentBookmark
BEGIN
    UPDATE Comment SET bookmark_count = bookmark_count + 1 WHERE id = NEW.comment_id;
END;

CREATE TRIGGER decrement_comment_bookmark_count
AFTER DELETE ON CommentBookmark
BEGIN
    UPDATE Comment SET bookmark_count = bookmark_count - 1 WHERE id = NEW.comment_id;
END;



CREATE INDEX idx_post_user_id ON Post(user_id);
CREATE INDEX idx_comment_post_id ON Comment(post_id);
CREATE INDEX idx_comment_parent_id ON Comment(parent_comment_id);
CREATE INDEX idx_userfollows_follower ON UserFollows(follower_id);
CREATE INDEX idx_userfollows_followee ON UserFollows(followee_id);
CREATE INDEX idx_notifications_receiver ON Notification(receiver_id, is_read, created_at DESC);
