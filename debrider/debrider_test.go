package debrider_test

import (
	"github.com/davidderus/christopher/debrider"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestDebrider(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Debrider Suite")
}

var _ = Describe("NewDebrider", func() {
	Context("With a valid debrider", func() {
		It("Should return an instantiated debrider", func() {
			debriderInstance, debriderError := debrider.NewDebrider("AllDebrid", nil)

			Expect(debriderError).NotTo(HaveOccurred())

			By("Initializing the debrider instance")
			allDebridInstance := &debrider.AllDebrid{}
			allDebridInstance.Init()

			Expect(debriderInstance).To(BeEquivalentTo(allDebridInstance))
		})
	})

	Context("With an invalid debrider", func() {
		It("Should return an error", func() {
			_, debriderError := debrider.NewDebrider("Fake", nil)

			Expect(debriderError.Error()).To(Equal("Invalid debrider given"))
		})
	})
})
