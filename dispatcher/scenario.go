package dispatcher

import (
	"errors"
	"fmt"
)

// Scenario handle URIs accross the application
type Scenario struct {
	currentStep *Step
	runError    error
	startFunc   func()
	endFunc     func()
	steps       []*Step
}

// Play runs a scenario through all its steps
func (s *Scenario) Play(event *Event) {
	if s.startFunc != nil {
		s.startFunc()
	}

	var currentStep, nextStep *Step
	var runError error

	currentStep = s.CurrentStep()
	// No initial Step set
	if currentStep == nil {
		s.runError = errors.New("No initial step provided")
		return
	}

	for {
		runError = currentStep.Run(event)
		if runError != nil {
			s.runError = runError
			return
		}

		// Going to next step
		nextStep = s.NextStep()

		if nextStep != nil {
			// Updating step
			s.SetCurrentStep(nextStep)
			currentStep = nextStep
		} else {
			// No next step, breaking out
			if s.endFunc != nil {
				s.endFunc()
			}
			return
		}
	}
}

// OnStart defines a callback to run when the scenario starts
func (s *Scenario) OnStart(callback func()) *Scenario {
	s.startFunc = callback
	return s
}

// OnEnd defines a callback to run when the scenario ends
func (s *Scenario) OnEnd(callback func()) *Scenario {
	s.endFunc = callback
	return s
}

// SetInitialStep defines the first step for the scenario
func (s *Scenario) SetInitialStep(step string) error {
	foundStep := s.findStepByName(step)
	if foundStep == nil {
		return fmt.Errorf("Undefined initial step %s", step)
	}

	s.SetCurrentStep(foundStep)

	return nil
}

// SetCurrentStep updates the scenario current step
func (s *Scenario) SetCurrentStep(step *Step) {
	s.currentStep = step
}

// CurrentStep returns the current scenario step
func (s *Scenario) CurrentStep() *Step {
	return s.currentStep
}

// RunError return a run error if any
func (s *Scenario) RunError() error {
	return s.runError
}

// From defines a step with the given name for the current scenario
func (s *Scenario) From(name string) *Step {
	newStep := &Step{from: name}
	s.steps = append(s.steps, newStep)

	return newStep
}

// NextStep returns next step to run, based on the current step next info
func (s *Scenario) NextStep() *Step {
	nextStepName := s.currentStep.Next()

	if nextStepName != "" {
		return s.findStepByName(nextStepName)
	}

	return nil
}

// findStepByName returns a step if it exists
func (s *Scenario) findStepByName(name string) *Step {
	for stepIndex, step := range s.steps {
		if step.From() == name {
			return s.steps[stepIndex]
		}
	}

	return nil
}
