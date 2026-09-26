package repositories

import (
	"database/sql"
	"strings"
	"time"

	"zone/internal/database"
	"zone/internal/followers/models"
)

type newFollower struct {
	id       int
	fullName string
}

func GetFollowStatus(requesterID, followeeID int) (string, error) {
	var status string

	err := database.Database.QueryRow(`SELECT status FROM followers WHERE follower_id = ? AND followee_id = ?`, requesterID, followeeID).Scan(&status)
	if err != nil {
		return "", err
	}

	return status, nil
}

func GetFollowCounts(userID int) (int, int, error) {
	var followersCount, followingCount int

	err := database.Database.QueryRow(`
		SELECT
			COALESCE(SUM(followee_id = ?), 0),
			COALESCE(SUM(follower_id = ?), 0)
		FROM followers
		WHERE status = 'accepted' AND (followee_id = ? OR follower_id = ?)`,
		userID, userID, userID, userID,
	).Scan(&followersCount, &followingCount)
	if err != nil {
		return 0, 0, err
	}

	return followersCount, followingCount, nil
}

func IsPrivateProfile(profileID int) (bool, error) {
	var isPrivate bool

	err := database.Database.QueryRow(
		"SELECT is_private FROM users WHERE id = ?", profileID,
	).Scan(&isPrivate)
	if err != nil {
		return false, err
	}

	return isPrivate, nil
}

func IsAcceptedFollower(viewerID, profileID int) (bool, error) {
	var isFollower bool

	err := database.Database.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM followers WHERE follower_id = ? AND followee_id = ? AND status = 'accepted')",
		viewerID, profileID,
	).Scan(&isFollower)
	if err != nil {
		return false, err
	}

	return isFollower, nil
}

func ListConnections(profileID int, listFollowers bool) ([]models.UserItem, error) {
	query := `
        SELECT u.id, u.first_name, u.last_name, COALESCE(u.nickname, ''), COALESCE(u.avatar, '')
        FROM followers f
        JOIN users u ON u.id = f.follower_id
        WHERE f.followee_id = ? AND f.status = 'accepted'`

	if !listFollowers {
		query = `
            SELECT u.id, u.first_name, u.last_name, COALESCE(u.nickname, ''), COALESCE(u.avatar, '')
            FROM followers f
            JOIN users u ON u.id = f.followee_id
            WHERE f.follower_id = ? AND f.status = 'accepted'`
	}

	rows, err := database.Database.Query(query, profileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]models.UserItem, 0)
	for rows.Next() {
		var u models.UserItem
		if err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Nickname, &u.Avatar); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func GetFollowTargetState(tx *sql.Tx, followerID, followeeID int) (isPrivate bool, hasStatus bool, status string, err error) {
	var rawStatus sql.NullString

	err = tx.QueryRow(`
        SELECT COALESCE(u.is_private, 0), f.status
        FROM users u
        LEFT JOIN followers f
            ON f.followee_id = u.id AND f.follower_id = ?
        WHERE u.id = ?`,
		followerID, followeeID,
	).Scan(&isPrivate, &rawStatus)
	if err != nil {
		return false, false, "", err
	}

	return isPrivate, rawStatus.Valid, rawStatus.String, nil
}

func CreateFollow(tx *sql.Tx, followerID, followeeID int, status string) error {
	_, err := tx.Exec(
		"INSERT INTO followers (follower_id, followee_id, status) VALUES (?, ?, ?)",
		followerID, followeeID, status,
	)

	return err
}

func DeleteFollow(tx *sql.Tx, followerID, followeeID int) error {
	_, err := tx.Exec(
		"DELETE FROM followers WHERE follower_id = ? AND followee_id = ?",
		followerID, followeeID,
	)

	return err
}

func AcceptPendingFollow(tx *sql.Tx, followerID, followeeID int) error {
	_, err := tx.Exec(
		"UPDATE followers SET status = 'accepted' WHERE follower_id = ? AND followee_id = ? AND status = 'pending'",
		followerID, followeeID,
	)

	return err
}

func AcceptFollowRequest(tx *sql.Tx, followerID, followeeID int) (int64, error) {
	result, err := tx.Exec(
		"UPDATE followers SET status = 'accepted' WHERE follower_id = ? AND followee_id = ? AND status = 'pending'",
		followerID, followeeID,
	)
	if err != nil {
		return 0, err
	}

	rowsAffected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		return 0, nil
	}

	return rowsAffected, nil
}

func DeclineFollowRequest(tx *sql.Tx, followerID, followeeID int) (int64, error) {
	result, err := tx.Exec(
		"DELETE FROM followers WHERE follower_id = ? AND followee_id = ? AND status = 'pending'",
		followerID, followeeID,
	)
	if err != nil {
		return 0, err
	}

	rowsAffected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		return 0, nil
	}

	return rowsAffected, nil
}

func AcceptPendingFollowers(tx *sql.Tx, userID int) error {
	rows, err := tx.Query(`
		SELECT f.follower_id, u.full_name
		FROM followers f
		JOIN users u ON u.id = f.follower_id
		WHERE f.followee_id = ? AND f.status = 'pending'`, userID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var followers []newFollower
	for rows.Next() {
		var nf newFollower
		if err := rows.Scan(&nf.id, &nf.fullName); err != nil {
			return err
		}
		followers = append(followers, nf)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	if _, err := tx.Exec(
		"UPDATE followers SET status = 'accepted' WHERE followee_id = ? AND status = 'pending'",
		userID,
	); err != nil {
		return err
	}

	if _, err := tx.Exec(
		"DELETE FROM notifications WHERE user_id = ? AND type = 'follow_request'",
		userID,
	); err != nil {
		return err
	}

	for _, f := range followers {
		name := strings.TrimSpace(f.fullName)
		if name == "" {
			name = "Someone"
		}
		if _, err := tx.Exec(`
			INSERT INTO notifications (user_id, actor_id, type, ref_id, message, is_read, created_at)
			VALUES (?, ?, 'new_follower', ?, ?, 0, ?)`,
			userID, f.id, f.id, name+" started following you", time.Now().UnixMilli(),
		); err != nil {
			return err
		}
	}

	return nil
}
