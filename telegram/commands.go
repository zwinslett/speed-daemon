package telegram

type Command = string

const (
	CmdLatest         Command = "/latest"
	CmdWeekly         Command = "/weekly"
	CmdMonthly        Command = "/monthly"
	CmdWeeklyCompare  Command = "/weekly_compare"
	CmdMonthlyCompare Command = "/monthly_compare"
)
