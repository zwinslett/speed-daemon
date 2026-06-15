package cmd

type Schedule = string

const (
	DailySchedule   Schedule = "0 8 1 * *"
	WeeklySchedule  Schedule = "0 8 * * 0"
	MonthlySchedule Schedule = "0 8 * * *"
)
