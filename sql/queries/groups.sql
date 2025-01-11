-- name: GetGroupIDByGroupName :one
SELECT id FROM groups WHERE group_name = $1;
