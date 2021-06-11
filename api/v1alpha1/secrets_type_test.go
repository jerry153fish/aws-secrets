package v1alpha1

import (
	"encoding/json"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSecretType(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Secret Type Suite")
}

var _ = Describe("CloudformationOutput", func() {
	var (
		cfn CloudformationOutput
	)
	BeforeEach(func() {
		cfn = CloudformationOutput{
			StackName: "syd-dev-test",
			KeyName:   "dbPass",
			OutputKey: "TestDBPassword",
		}

	})

	Describe("CloudformationOutput marshal", func() {
		Context("Simple test without type", func() {
			It("Should match predefined result", func() {
				result, _ := json.Marshal(cfn)
				Expect(string(result)).To(Equal(`{"stackName":"syd-dev-test","key":"dbPass","outputKey":"TestDBPassword"}`))
			})
		})

		// Context("With fewer than 300 pages", func() {
		//     It("should be a short story", func() {
		//         Expect(shortBook.CategoryByLength()).To(Equal("SHORT STORY"))
		//     })
		// })
	})
})
