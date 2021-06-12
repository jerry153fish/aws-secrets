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
		cfn       CloudformationOutput
		cfnString string
	)
	BeforeEach(func() {
		cfn = CloudformationOutput{
			StackName: "syd-dev-test",
			KeyName:   "dbPass",
			OutputKey: "TestDBPassword",
		}

		cfnString = `{"stackName":"syd-dev-test","key":"dbPass","outputKey":"TestDBPassword"}`

	})

	Describe("CloudformationOutput marshal/unmarshal", func() {
		Context("Marsh Cloudformation", func() {
			It("Should match string", func() {
				result, _ := json.Marshal(cfn)

				Expect(string(result)).To(Equal(cfnString))
			})
		})

		Context("Unmarsh Cloudformation", func() {
			It("should match fields", func() {
				var result CloudformationOutput
				if err := json.Unmarshal([]byte(cfnString), &result); err != nil {
					panic(err)
				}
				Expect(result.StackName).To(Equal("syd-dev-test"))
				Expect(result.KeyName).To(Equal("dbPass"))
				Expect(result.OutputKey).To(Equal("TestDBPassword"))
			})
		})
	})
})
