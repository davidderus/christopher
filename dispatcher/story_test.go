package dispatcher_test

import (
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/davidderus/christopher/dispatcher"
)

type TestStory struct {
	logChan chan string
}

func (ts *TestStory) SetLog(logChannel chan string) {
	ts.logChan = logChannel
}

func (ts *TestStory) Scenario() *Scenario {
	scenario := &Scenario{}

	scenario.OnStart(func() {
		ts.logChan <- "Starting story"
	}).OnEnd(func() {
		ts.logChan <- "Ending story"
	})

	scenario.From("submitter").To("transformer").Do(func(event *Event) error {
		ts.logChan <- fmt.Sprintf("%s submitted", event.Value)

		if event.Value != "http://google.fr" {
			return errors.New("Not the expected URL")
		}

		return nil
	}).OnStart(func() {
		ts.logChan <- "start submitter"
	}).OnEnd(func() {
		ts.logChan <- "end submitter"
	})

	scenario.From("transformer").Do(func(event *Event) error {
		// Updating event origin as we act on it
		event.Origin = "transformer"
		oldValue := event.Value
		event.Value = oldValue + "/new"

		ts.logChan <- fmt.Sprintf("%s is transformed to %s", oldValue, event.Value)

		return nil
	})

	return scenario
}

// NewStory returns a Christopher story (by default a BasicStory)
func NewTestStory(logChan chan string) Story {
	story := &TestStory{}
	story.SetLog(logChan)

	return story
}

var _ = Describe("Story", func() {
	Context("With a custom story", func() {
		It("should work", func(done Done) {
			channel := make(chan string, 0)
			defer close(channel)

			story := NewTestStory(channel)
			scenario := story.Scenario()

			By("Setting an initial step")
			initialStepError := scenario.SetInitialStep("submitter")

			Expect(initialStepError).NotTo(HaveOccurred())
			Expect(scenario.CurrentStep().From()).To(Equal("submitter"))

			By("Playing an event")
			customEvent := &Event{Origin: "cli", Value: "http://google.fr"}
			go scenario.Play(customEvent)

			Expect(<-channel).To(Equal("Starting story"))
			Expect(<-channel).To(Equal("start submitter"))
			Expect(<-channel).To(Equal("http://google.fr submitted"))
			Expect(<-channel).To(Equal("end submitter"))
			Expect(<-channel).To(Equal("http://google.fr is transformed to http://google.fr/new"))
			Expect(<-channel).To(Equal("Ending story"))

			Expect(scenario.RunError()).To(BeNil())
			Expect(scenario.CurrentStep().From()).To(Equal("transformer"))

			Expect(customEvent.Origin).To(Equal("transformer"))
			Expect(customEvent.Value).To(Equal("http://google.fr/new"))

			close(done)
		})

		It("should raise an error if something is wrong on Do", func(done Done) {
			channel := make(chan string, 0)
			defer close(channel)

			story := NewTestStory(channel)
			scenario := story.Scenario()

			By("Setting an initial step")
			initialStepError := scenario.SetInitialStep("submitter")

			Expect(initialStepError).NotTo(HaveOccurred())
			Expect(scenario.CurrentStep().From()).To(Equal("submitter"))

			By("Playing an event")
			customEvent := &Event{Origin: "cli", Value: "http://google.com"}
			go scenario.Play(customEvent)

			Expect(<-channel).To(Equal("Starting story"))
			Expect(<-channel).To(Equal("start submitter"))
			Expect(<-channel).To(Equal("http://google.com submitted"))

			Expect(scenario.RunError()).NotTo(BeNil())
			Expect(scenario.RunError().Error()).To(Equal("Not the expected URL"))

			close(done)
		})
	})
})
