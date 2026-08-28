import { browser } from '$app/environment';
import { api } from '$lib/api/client';
import { pushToast } from '$lib/stores/toast';

export async function requestReminderPermission(): Promise<NotificationPermission | 'unsupported'> {
	if (!browser || !('Notification' in window)) return 'unsupported';
	return Notification.requestPermission();
}
export async function evaluateReminders(): Promise<void> {
	try {
		const reminders = await api.post<Array<{ id: string; state: string; questionId?: string }>>(
			'/api/v1/reminders/evaluate'
		);
		for (const reminder of reminders) {
			if (reminder.state === 'delivered' || reminder.state === 'missed')
				pushToast(
					reminder.state === 'missed'
						? 'A reminder was missed while Noted was unavailable.'
						: 'A question reminder is due.',
					'info',
					8000
				);
		}
	} catch {
		/* The local service may be restarting; the next resume retries. */
	}
}
