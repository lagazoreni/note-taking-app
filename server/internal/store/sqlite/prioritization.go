package sqlite

import (
	"context"
	"database/sql"
	"time"

	"noted.local/noted/internal/domain"
	"noted.local/noted/internal/platform"
)

func (s *CaptureStore) SetReminder(ctx context.Context, questionID string, scheduledAt time.Time) (domain.Reminder, error) {
	question, err := s.GetQuestion(ctx, questionID)
	if err != nil {
		return domain.Reminder{}, err
	}
	now := s.now()
	existing, err := s.reminder(ctx, questionID)
	if err != nil {
		return domain.Reminder{}, err
	}
	if existing == nil {
		id, e := s.newID()
		if e != nil {
			return domain.Reminder{}, e
		}
		_, err = s.db.ExecContext(ctx, "INSERT INTO reminders(id,question_id,scheduled_at,state,last_evaluated_at,created_at,updated_at,version) VALUES(?,?,?,'pending',NULL,?,?,1)", id, questionID, platform.FormatTimestamp(scheduledAt), now, now)
	} else {
		_, err = s.db.ExecContext(ctx, "UPDATE reminders SET scheduled_at=?,state='pending',last_evaluated_at=NULL,updated_at=?,version=version+1 WHERE question_id=?", platform.FormatTimestamp(scheduledAt), now, questionID)
	}
	if err != nil {
		return domain.Reminder{}, err
	}
	value, err := s.reminder(ctx, questionID)
	if err != nil || value == nil {
		return domain.Reminder{}, err
	}
	_ = question
	return *value, nil
}
func (s *CaptureStore) ClearReminder(ctx context.Context, questionID string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM reminders WHERE question_id=?", questionID)
	return err
}
func (s *CaptureStore) EvaluateReminders(ctx context.Context, now time.Time) ([]domain.Reminder, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id,question_id,scheduled_at,state,last_evaluated_at,created_at,updated_at,version FROM reminders WHERE state='pending' AND scheduled_at <= ? ORDER BY scheduled_at,id", platform.FormatTimestamp(now))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]domain.Reminder, 0)
	for rows.Next() {
		var r domain.Reminder
		var questionID, scheduled, created, updated string
		var last sql.NullString
		if err := rows.Scan(&r.ID, &questionID, &scheduled, &r.State, &last, &created, &updated, &r.Version); err != nil {
			return nil, err
		}
		r.ScheduledAt, _ = platform.ParseTimestamp(scheduled)
		r.CreatedAt, _ = platform.ParseTimestamp(created)
		r.UpdatedAt, _ = platform.ParseTimestamp(updated)
		if last.Valid {
			value, _ := platform.ParseTimestamp(last.String)
			r.LastEvaluatedAt = &value
		}
		newState := domain.ReminderDelivered
		if now.Sub(r.ScheduledAt) > 24*time.Hour {
			newState = domain.ReminderMissed
		}
		if _, err := s.db.ExecContext(ctx, "UPDATE reminders SET state=?,last_evaluated_at=?,updated_at=?,version=version+1 WHERE id=? AND version=?", newState, platform.FormatTimestamp(now), platform.FormatTimestamp(now), r.ID, r.Version); err != nil {
			return nil, err
		}
		r.State = newState
		r.Version++
		result = append(result, r)
	}
	return result, rows.Err()
}
func (s *CaptureStore) DismissReminder(ctx context.Context, id string, version int64) error {
	result, err := s.db.ExecContext(ctx, "UPDATE reminders SET state='dismissed',updated_at=?,version=version+1 WHERE id=? AND version=? AND state IN ('delivered','missed')", s.now(), id, version)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return ErrConflict
	}
	return nil
}
