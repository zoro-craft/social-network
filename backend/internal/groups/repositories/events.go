package repositories

import (
	"zone/internal/database"
	"zone/internal/groups/models"
)

func CreateEvent(event models.Event) (int64, error) {
	result, err := database.Database.Exec(
		"INSERT INTO events (group_id, user_id, title, description, event_time) VALUES (?, ?, ?, ?, ?)",
		event.GroupID, event.UserID, event.Title, event.Description, event.EventTime,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func GetEvent(eventID int) (models.Event, error) {
	e := models.Event{}
	err := database.Database.QueryRow(`
		SELECT ev.id, ev.group_id, ev.user_id, ev.title, COALESCE(ev.description, ''),
		       ev.event_time, ev.created_at
		FROM events ev
		WHERE ev.id = ?`, eventID,
	).Scan(&e.ID, &e.GroupID, &e.UserID, &e.Title, &e.Description, &e.EventTime, &e.CreatedAt)
	return e, err
}

func ListEvents(groupID int, viewerID int) ([]models.Event, error) {
	rows, err := database.Database.Query(`
		SELECT ev.id, ev.group_id, ev.user_id, ev.title, COALESCE(ev.description, ''),
		       ev.event_time, ev.created_at, COALESCE(u.first_name, ''),
		       (SELECT COUNT(*) FROM event_responses er WHERE er.event_id = ev.id AND er.response = 'going'),
		       (SELECT COUNT(*) FROM event_responses er WHERE er.event_id = ev.id AND er.response = 'not_going'),
		       (SELECT COUNT(*) FROM event_responses er WHERE er.event_id = ev.id AND er.response = 'maybe'),
		       COALESCE((SELECT er.response FROM event_responses er WHERE er.event_id = ev.id AND er.user_id = ?), '')
		FROM events ev
		JOIN users u ON ev.user_id = u.id
		WHERE ev.group_id = ?
		ORDER BY ev.event_time ASC`, viewerID, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []models.Event{}
	for rows.Next() {
		var e models.Event
		if err := rows.Scan(&e.ID, &e.GroupID, &e.UserID, &e.Title, &e.Description, &e.EventTime, &e.CreatedAt, &e.CreatorName, &e.Going, &e.NotGoing, &e.Maybe, &e.MyResponse); err == nil {
			events = append(events, e)
		}
	}
	return events, rows.Err()
}

func SetEventResponse(eventID, userID int, response string) error {
	_, err := database.Database.Exec(`
		INSERT INTO event_responses (event_id, user_id, response) VALUES (?, ?, ?)
		ON CONFLICT(event_id, user_id) DO UPDATE SET response = excluded.response`,
		eventID, userID, response,
	)
	return err
}
