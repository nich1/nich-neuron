package main

const calendarRemindersFile = "database/time/calendar/reminders.json"

func createReminder(event calendarEvent) {

}

func deleteReminder(event reminder) {

}

type eventAndID struct {
	event calendarEvent
	ID    uint8
}

func deleteReminderByEventAndID(eventStr eventAndID)
