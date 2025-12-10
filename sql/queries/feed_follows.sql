-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING *
)
SELECT
    inserted_feed_follow.*,
    feeds.name as feedName,
    users.name as userName
FROM inserted_feed_follow
INNER JOIN feeds
ON feeds.id = inserted_feed_follow.feed_id
INNER JOIN users
ON users.id = inserted_feed_follow.user_id;


-- name: GetFeedFollowsForUser :many
SELECT * FROM feed_follows
WHERE user_id = $1;

-- name: GetFeedFollowByFeedId :one
SELECT *, feeds.name as feedName
FROM feed_follows
INNER JOIN feeds
ON feed_follows.feed_id = feeds.id
WHERE feed_follows.feed_id = $1;

-- name: DeleteByUserIdAndFeedId :execrows
DELETE FROM feed_follows
WHERE user_id = $1
AND feed_id = $2;
