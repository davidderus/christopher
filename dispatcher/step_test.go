package dispatcher_test

import (
	. "github.com/davidderus/christopher/dispatcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Step", func() {
	var (
		basicEvent      *Event
		currentScenario *Scenario
	)

	BeforeEach(func() {
		currentScenario = &Scenario{}
		basicEvent = &Event{Origin: "test", Value: "my-value"}
	})

	Describe(".If()", func() {
		Context("without the If() condition", func() {
			It("should execute step2", func() {
				var isValid bool

				baseString := "Hello"

				currentScenario.From("step1").To("step2").Do(func(_ *Event) error {
					isValid = true

					return nil
				})

				currentScenario.From("step2").Do(func(_ *Event) error {
					baseString = "World"

					return nil
				})

				currentScenario.SetInitialStep("step1")
				currentScenario.Play(basicEvent)

				Expect(baseString).To(Equal("World"))
			})
		})

		Context("With the If() condition", func() {
			It("should execute step if condition is fullfilled", func() {
				var isValid bool

				baseString := "Hello"

				currentScenario.From("step1").To("step2").Do(func(_ *Event) error {
					isValid = true

					return nil
				})

				currentScenario.From("step2").Do(func(_ *Event) error {
					baseString = "World"

					return nil
				}).If(func() bool { return isValid })

				currentScenario.SetInitialStep("step1")
				currentScenario.Play(basicEvent)

				// Step is run, as if isValid was evaluated at instanciation
				Expect(baseString).To(Equal("World"))
			})

			It("should skip step if condition is false", func() {
				isValid := true
				baseString := "Hello"

				currentScenario.From("step1").To("step2").Do(func(_ *Event) error {
					isValid = false

					return nil
				})

				currentScenario.From("step2").Do(func(_ *Event) error {
					baseString = "World"

					return nil
				}).If(func() bool { return isValid })

				currentScenario.SetInitialStep("step1")
				currentScenario.Play(basicEvent)

				// Step is run, as if isValid was evaluated at instanciation
				Expect(baseString).To(Equal("Hello"))
			})
		})
	})
})
