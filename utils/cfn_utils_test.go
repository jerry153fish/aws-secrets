package utils

import (
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickmn/go-cache"
)

func TestSecretType(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cloudformation Utils Suite")
}

var _ = Describe("Cloudformation Utils", func() {
	var (
		cf *cloudformation.CloudFormation
		c  *cache.Cache
	)
	BeforeEach(func() {
		os.Setenv("TEST_ENV", "YES")
		cf = GetCfnClient()
		c = cache.New(5*time.Minute, 10*time.Minute)
	})

	Describe("Cfn until", func() {
		Context("GetStackOutput", func() {
			It("Should match fake s3 bucket output", func() {

				result, err := GetStackOutput(cf, "myteststack", "S3Bucket", c)
				Expect(err).To(BeNil())
				Expect(result).To(Equal("S3Bucket"))
			})
		})
	})
})
