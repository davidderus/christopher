package dispatcher

// Event represents an event going through the Story
type Event struct {
	Value  string // A valid URI
	Origin string // Previous handler (submitter, debrider, downloader…)
}

// Story is the implementation of a scenario
type Story interface {
	// Scenario defines the story scenario
	Scenario() *Scenario
}
