-- name: CreateFeed :one
INSERT INTO
  feeds (id, created_at, updated_at, name, url, user_id)
VALUES
  ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: GetAllFeeds :many
SELECT
  *
FROM
  feeds;

-- name: GetFeedByUrl :one
SELECT
  *
FROM
  feeds
WHERE
  url = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET
  updated_at = Now(),
  last_fetched_at = Now()
WHERE
  feeds.id = $1;

-- name: GetNextFeedToFetch :one
SELECT
  *
FROM
  feeds
ORDER BY
  last_fetched_at NULLS FIRST
LIMIT
  1;
