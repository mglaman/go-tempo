package tempo

type WorkAttributeValue struct {
	Key   string
	Value string
}

type Worklog struct {
	Self           string
	TempoWorklogId int
	JiraWorklogId  int
	Issue          struct {
		Self string
		Key  string
	}
	// Set to float so we can perform math
	TimeSpentSeconds float32
	StartDate        string
	StartTime        string
	Description      string
	CreatedAt        string
	UpdatedAt        string
	Author           User
	Attributes       struct {
		Self  string
		Items []WorkAttributeValue
	}
}
type WorklogCollection struct {
	Self     string
	Metadata struct {
		Count    int
		Offset   int
		Limit    int
		Next     string
		Previous string
	}
	Results []Worklog
}

type WorklogPayload struct {
	IssueKey         string  `json:"issueKey"`
	TimeSpentSeconds float64 `json:"timeSpentSeconds"`
	BillableSeconds  float64 `json:"billableSeconds"`
	StartDate        string  `json:"startDate"`
	StartTime        string  `json:"startTime"`
	Description      string  `json:"description"`
	AuthorUsername   string  `json:"authorUsername"`
}
