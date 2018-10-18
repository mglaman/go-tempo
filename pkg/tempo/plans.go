package tempo

type PlanPeriod struct {
	From string
	To   string
	// Define as a float so we can convert it to hours and allow a decimal.
	TimePlannedSeconds float32
}

type Plan struct {
	Self        string
	Id          int
	StartDate   string
	EndDate     string
	CreatedAt   string
	Description string
	UpdatedAt   string
	Assignee    struct {
		Self string
		Type string
	}
	PlanItem struct {
		Self string
		Type string
	}
	Recurrence struct {
		Rule              string
		RecurrenceEndDate string
	}
	Dates struct {
		Metadata struct {
			Count int
			All   string
		}
		Values []PlanPeriod
	}
}
type PlanCollection struct {
	Self     string
	Metadata struct {
		Count int
	}
	Results []Plan
}
