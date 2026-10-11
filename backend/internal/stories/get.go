package stories

import (
	"net/http"
	"time"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/utils"
)

// StoryLifetime is how long a story stays visible, matching Instagram's 24h window.
const StoryLifetime = 24 * time.Hour

const listStoriesQuery = `
	SELECT
		s.id,
		s.user_id,
		u.full_name,
		COALESCE(u.avatar, ''),
		s.image,
		s.caption,
		s.created_at,
		CASE WHEN sv.story_id IS NULL THEN 0 ELSE 1 END AS seen
	FROM stories s
	JOIN users u ON s.user_id = u.id
	LEFT JOIN story_views sv ON sv.story_id = s.id AND sv.user_id = ?
	WHERE s.created_at > ?
	  AND (
		s.user_id = ?
		OR EXISTS (
			SELECT 1 FROM followers f
			WHERE f.follower_id = ? AND f.followee_id = s.user_id AND f.status = 'accepted'
		)
	  )
	ORDER BY s.created_at ASC
`

// storyVisibleCondition reuses the same audience rules as listStoriesQuery:
// the author always sees their own stories, everyone else must be an accepted follower.
const storyVisibleCondition = `
	s.id = ? AND s.created_at > ? AND (
		s.user_id = ?
		OR EXISTS (
			SELECT 1 FROM followers f
			WHERE f.follower_id = ? AND f.followee_id = s.user_id AND f.status = 'accepted'
		)
	)
`

func GetStoriesAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	viewerID, err := auth.GetUserID(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "You must be logged in")
		return
	}

	groups, err := GetActiveStories(viewerID)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to load stories")
		return
	}

	utils.RespondJSON(w, http.StatusOK, groups)
}

// GetActiveStories returns unexpired stories grouped by author, oldest first so the
// client can play a group in the order it was posted. Only the viewer's own stories
// and those of people they follow are included.
func GetActiveStories(viewerID int) ([]StoryGroup, error) {
	cutoff := time.Now().Add(-StoryLifetime).UnixMilli()

	rows, err := database.Database.Query(listStoriesQuery, viewerID, cutoff, viewerID, viewerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// A group starts out "seen" and flips to unseen as soon as one story is new.
	groups := []StoryGroup{}
	index := map[int]int{}

	for rows.Next() {
		var story Story
		var fullName, avatar string
		var seen bool

		if err := rows.Scan(
			&story.ID, &story.UserID, &fullName, &avatar,
			&story.Image, &story.Caption, &story.CreatedAt, &seen,
		); err != nil {
			return nil, err
		}

		story.Seen = seen

		groupIndex, ok := index[story.UserID]
		if !ok {
			groups = append(groups, StoryGroup{
				UserID:   story.UserID,
				FullName: fullName,
				Avatar:   avatar,
				Seen:     true,
				Stories:  []Story{},
			})
			groupIndex = len(groups) - 1
			index[story.UserID] = groupIndex
		}

		groups[groupIndex].Stories = append(groups[groupIndex].Stories, story)
		if !seen {
			groups[groupIndex].Seen = false
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return groups, nil
}

// CanViewStory reports whether the viewer may see an unexpired story.
func CanViewStory(storyID, viewerID int) bool {
	cutoff := time.Now().Add(-StoryLifetime).UnixMilli()

	var exists bool
	database.Database.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM stories s WHERE "+storyVisibleCondition+")",
		storyID, cutoff, viewerID, viewerID,
	).Scan(&exists)
	return exists
}
