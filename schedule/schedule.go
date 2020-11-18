package schedule

// Season represents one season of school where the schedule stays same
type Season struct {
	Name  string
	Start int64
	End   int64
}

// Schedule represents some users schedule
type Schedule struct {
	Groups  []Group
	OwnerID string
	Season  Season
}

// AddSchedule creates an schedule for user in season and saves it to database
func AddSchedule() (Schedule, error) {
	// TODO Create and save schedule
	return Schedule{}, nil
}

// GetSchedule gets the schedule of an user in specific season
func GetSchedule(ownerID string, season int) (Schedule, error) {
	// TODO Cet the schedule of user for specific season
	return Schedule{}, nil
}
