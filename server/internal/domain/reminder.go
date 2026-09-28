package domain

import "fmt"

const (
	ReminderPending   = "pending"
	ReminderDelivered = "delivered"
	ReminderMissed    = "missed"
	ReminderDismissed = "dismissed"
)

func ValidReminderState(state string) bool {
	return state == ReminderPending || state == ReminderDelivered || state == ReminderMissed || state == ReminderDismissed
}
func ValidateReminderState(state string) error {
	if !ValidReminderState(state) {
		return fmt.Errorf("invalid reminder state")
	}
	return nil
}
