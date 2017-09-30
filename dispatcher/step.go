package dispatcher

import "fmt"

// Step is a step during a Scenario
type Step struct {
	from string
	to   string

	conditionFunc func() bool

	doFunc      func(event *Event) error
	onStartFunc func()
	onEndFunc   func()
}

// Do defines something to do during step
func (s *Step) Do(doFunc func(event *Event) error) *Step {
	s.doFunc = doFunc
	return s
}

// OnStart defines something done before the step is executed
func (s *Step) OnStart(onStartFunc func()) *Step {
	s.onStartFunc = onStartFunc
	return s
}

// OnEnd defines something done after the step is successfully executed
func (s *Step) OnEnd(onEndFunc func()) *Step {
	s.onEndFunc = onEndFunc
	return s
}

// Run runs the step and its callbacks
func (s *Step) Run(event *Event) error {
	// Skipping if condition is false
	if s.conditionFunc != nil {
		currentConditionState := s.conditionFunc()
		if !currentConditionState {
			return nil
		}
	}

	if s.onStartFunc != nil {
		s.onStartFunc()
	}

	if s.doFunc != nil {
		doError := s.doFunc(event)
		if doError != nil {
			return doError
		}
	} else {
		return fmt.Errorf("Nothing to do in step %s", s.from)
	}

	if s.onEndFunc != nil {
		s.onEndFunc()
	}

	return nil
}

// To defines the step next to the current one
func (s *Step) To(nextStep string) *Step {
	s.to = nextStep
	return s
}

// Next returns the step name next to the current step
func (s *Step) Next() string {
	return s.to
}

// From returns the current step from()
func (s *Step) From() string {
	return s.from
}

// If defines an execution condition for current step
//
// NOTE conditionFunc is a callback evaluated at runtime. You may want to wrap
// your step in a if condition if you don't need such a dynamic evaluation.
func (s *Step) If(conditionFunc func() bool) *Step {
	s.conditionFunc = conditionFunc
	return s
}
