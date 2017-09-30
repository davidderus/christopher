package dispatcher_test

import (
	. "github.com/davidderus/christopher/dispatcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Scenario", func() {
	Context("without initial step set", func() {
		It("should raise an error", func() {
			scenario := &Scenario{}
			scenario.From("tested")

			customEvent := &Event{Origin: "CLI", Value: "http://google.com"}
			scenario.Play(customEvent)

			Expect(scenario.RunError()).NotTo(BeNil())
			Expect(scenario.RunError().Error()).To(Equal("No initial step provided"))
		})
	})

	Context("with invalid Do Action", func() {
		It("should raise an error", func() {
			scenario := &Scenario{}
			scenario.From("tested")
			scenario.SetInitialStep("tested")

			customEvent := &Event{Origin: "CLI", Value: "http://google.com"}
			scenario.Play(customEvent)

			Expect(scenario.RunError()).NotTo(BeNil())
			Expect(scenario.RunError().Error()).To(Equal("Nothing to do in step tested"))
		})
	})

	Context("with an invalid initial step", func() {
		It("should return some errors", func() {
			scenario := &Scenario{}

			By("Setting an initial step")
			initialStepError := scenario.SetInitialStep("invalid")

			Expect(initialStepError.Error()).To(Equal("Undefined initial step invalid"))
		})
	})
})
